package models

import (
	"Goez/pkg/e"
	"Goez/pkg/search"
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"strconv"
)

// 文章
type Article struct {
	ID      int    `gorm:"primary_key" json:"id"`
	Title   string `gorm:"varchar(50);unique" json:"title"`
	Desc    string `gorm:"varchar(200)" json:"desc"`
	Content string `gorm:"type:text" json:"content"`
	View    int64  `gorm:"default:0" json:"view"`
	UserId  int    `gorm:"index" json:"user_id"`
	Tags    []Tag  `gorm:"many2many:article_tags;"`
	gorm.Model
}

// 文章列表
type ArticleList struct {
	ID    int
	Title string
	View  int
}

// 获取文章
func (art Article) GetArticleByFields(field, keyword string, pageIndex, pageSize int) (al []ArticleList, count int, err error) {

	fields := []string{"title", "desc", "content"}

	if field != "" {
		if ok, err := search.In(fields, field); !ok {
			return al, count, err
		}
	}

	if field == "" && keyword == "" {
		Repo.SqlClient.Debug().Scopes(Paginate(pageIndex, pageSize)).Model(&Article{}).Select("id,title,view").Scan(&al)
	} else {
		err := Repo.SqlClient.Debug().Model(&Article{}).
			Select("id,title,view").
			Scan(&al).
			Order("view desc").Error

		if err != nil {
			return al, count, err
		}
	}
	count = len(al)

	return al, count, nil
}

// 根据id 获取文章
func (art Article) GetArticlesById() (a Article, err error) {

	//// 查询Redis key
	a, err = CacheArticle(art.ID)

	//// 返回缓存
	if err == nil {

		//更新缓存访问量
		if err := CacheSetArticleView(a); err != nil {
			return a, err
		}

		return a, err
	}

	//Repo.SqlClient.Where("title like ?", "%"+art.Title+"%").First(&art)
	a.ID = art.ID

	if err := Repo.SqlClient.Model(&Article{}).Preload("Tags").First(&a).Error; err != nil {
		return a, err
	}

	// 设置缓存
	if err := a.CacheSetArticle(); err != nil {
		return a, err
	}

	// 访问量 + 1
	if err := a.ViewArticle(); err != nil {
		return a, err
	}

	return a, nil
}

// 添加文章
func (art Article) AddArticle() (Article, error) {

	ts := Repo.SqlClient.Begin()
	defer ts.Rollback()

	u := User{ID: art.UserId}

	err := ts.First(&u).Error

	if err != nil {
		return Article{}, err
	}

	// TODO CHECK ARTICLE PARAMS

	err = ts.Create(&art).Error

	if err != nil {
		return Article{}, err
	}

	return art, ts.Commit().Error
}

// 删除文章 TODO 清除关系
func (art Article) DeleteArticle() error {

	aa, err := art.GetArticlesById()

	if err := Repo.SqlClient.Model(&aa).Association("Tags").Delete(aa.Tags); err != nil {
		return err
	}

	if err := Repo.SqlClient.Debug().Delete(&aa).Error; err != nil {
		return err
	}

	return err
}

// 更新文章
func (art Article) UpdateArticle() error {
	return Repo.SqlClient.Save(&art).Error
}

// 获取文章列表
func (art Article) GetArticleTitleList() (al []ArticleList, err error) {

	al, err = CacheQueryArticleList()
	if err == nil {
		return al, err
	} else {
		// err Redis 连接异常等情况
	}

	err = Repo.SqlClient.Debug().Model(&Article{}).Select("id,title,view").Order("view desc").Limit(10).Scan(&al).Error

	if err != nil {
		return al, err
	}

	if err := CacheSetArticleList(al); err != nil {
		return al, err
	}

	return al, nil
}

// 查询所有文章
func ExportArticleData() (as []Article, err error) {

	err = Repo.SqlClient.Model(&Article{}).Preload("Tags").Find(&as).Error
	return
}

func GetArticleByIds(ids []int) (al []ArticleList, err error) {

	err = Repo.SqlClient.Debug().Model(&Article{}).Where("id in ?", ids).Scan(&al).Error

	return
}

// 更新
func (art Article) UpdateService() error {

	id := strconv.Itoa(art.ID)
	articleKey := e.CACHE_ARTICLE + id
	articleViewKey := e.CACHE_ARTICLE_VIEW + id

	if err := Repo.Redis.Del(articleKey).Err(); err != nil {
		return err
	}

	if err := art.UpdateArticle(); err != nil {
		return err
	}

	if view, err := Repo.Redis.Get(articleViewKey).Int64(); err != nil {
		return err
	} else {
		art.View = view
	}

	if err := Repo.SqlClient.First(&art).Error; err != nil {
		return err
	}

	if jsonStr, err := jsoniter.Marshal(art); err != nil {
		return err
	} else {
		return Repo.Redis.Set(articleKey, jsonStr, -1).Err()
	}
}

// ------------------------ Cache ------------------------

// 更新 文章访问db
func (a Article) ViewArticle() (err error) {

	if err = Repo.SqlClient.Debug().Model(&Article{}).Where("id = ?", a.ID).Update("view", gorm.Expr("view + ?", 1)).Error; err != nil {
		return err
	}
	return err
}

// 设置文章缓存
func (a Article) CacheSetArticle() (err error) {

	id := strconv.Itoa(a.ID)
	articleKey := e.CACHE_ARTICLE + id
	articleViewKey := e.CACHE_ARTICLE_VIEW + id

	jsonStr, err := jsoniter.Marshal(a)
	if err != nil {
		return err
	}

	if err := Repo.Redis.Set(articleKey, jsonStr, -1).Err(); err != nil {
		return err
	}

	if err := Repo.Redis.Set(articleViewKey, a.View, -1).Err(); err != nil {
		return err
	}

	return
}

// 查询文章缓存
func CacheArticle(articleId int) (a Article, err error) {

	id := strconv.Itoa(articleId)
	articleKey := e.CACHE_ARTICLE + id
	articleViewKey := e.CACHE_ARTICLE_VIEW + id

	adata, err := Repo.Redis.Get(articleKey).Result()
	view, err := Repo.Redis.Get(articleViewKey).Int64()

	// Redis 查询到key
	if err == nil {

		if err := json.Unmarshal([]byte(adata), &a); err != nil {
			return a, err
		}

		a.View = view

		return a, nil
	}
	return a, err
}

// 更新缓存文章访问量
func CacheSetArticleView(a Article) error {

	id := strconv.Itoa(a.ID)
	articleViewKey := e.CACHE_ARTICLE_VIEW + id
	articleKey := e.CACHE_ARTICLE + id
	//hasRecord := e.CACHE_USER_RECORD + id

	view := Repo.Redis.Incr(articleViewKey).Val()

	if view%1000 == 0 {
		if err := Repo.SqlClient.Model(&a).Update("view", view).Error; err != nil {
			return err
		}

		a.View = view

		jstr, _ := json.Marshal(a)

		if err := Repo.Redis.Set(articleKey, jstr, -1).Err(); err != nil {
			return err
		}
	}

	return nil
}

// 设置文章列表缓存
func CacheSetArticleList(al []ArticleList) error {
	return Repo.Redis.Set("current_article_list", al, 3600*2).Err()
}

func CacheQueryArticleList() (as []ArticleList, err error) {
	data, err := Repo.Redis.Get("current_article_list").Result()
	if err != nil {
		return as, err
	}

	if err := json.Unmarshal([]byte(data), &as); err != nil {
		return as, err
	}

	return as, nil
}
