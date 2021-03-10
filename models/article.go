package models

import (
	"Goez/pkg/config"
	"Goez/pkg/e"
	"Goez/pkg/gredis"
	"Goez/pkg/search"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"strconv"
)

type Article struct {
	ID          int 		`gorm:"primary_key" json:"id"`
	Title 		string 		`gorm:"varchar(50)" json:"title"`
	Desc		string 		`gorm:"varchar(200)" json:"desc"`
	Content     string 		`gorm:"type:text" json:"content"`
	View		int64		`gorm:"default:0" json:"view"`
	UserId		int			`gorm:"index" json:"user_id"`
	Tags        *[]Tag		`gorm:"many2many:article_tags;"`
	ModelTime
}


type ArticleList struct {
	ID				int
	Title			string
	View            int
}

// 获取文章
func (art Article) GetArticleByFields(field, keyword string) (al []ArticleList, count int, err error) {

	fields := []string{"title", "desc", "content"}

	if field != "" {
		if ok, err := search.In(fields, field); !ok {
			return al,count, err
		}
	}

	if field == "" && keyword == "" {
		db.Debug().Model(&Article{}).Select("id,title,view").Where("deleted_on is null").Scan(&al)
	} else {
		err := db.Debug().Model(&Article{}).
			Select("id,title,view").
			Where(field+" like  ? and deleted_on is null", "%"+keyword+"%").
			Scan(&al).
			Order("view desc").Error

		if err != nil {
		    return al,count, err
		}
	}
	count = len(al)


	return al, count, nil
}

func (art Article) GetArticlesById() (a Article, err error) {

	// 查询redis key 并 + 1访问量
	a, err = CacheArticle(art.ID)

	// 返回缓存
	if err == nil {
		// 更新缓存访问量
		if err := CacheSetArticleView(a); err != nil {
			return a, err
		}
		a.View = a.View + 1
	   return a, err
	}

	// 为查询到缓存， 查询mysql
	a.ID = art.ID

	if err := db.Debug().Model(&Article{}).First(&a).Error; err != nil {
		return Article{}, err
	}

	//// 访问量 + 1
	//if err := a.ViewArticle(); err != nil {
	//	return a, err
	//}

	// 设置缓存
	if err := a.CacheSetArticle(); err != nil {
		return a, err
	}

	return a, nil
}

func (art Article) AddArticle() (Article, error) {


	ts := db.Begin()
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

func (art Article)  DeleteArticle() (bool, error) {

	if err := db.Unscoped().Where("id = ?", art.ID).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (art Article)  GetArticleTitleList() (al []ArticleList, err error) {

	al, err = CacheQueryArticleList()
	if err == nil {
		return al, err
	} else {
		// err redis 连接异常等情况
	}

	err = db.Debug().Model(&Article{}).Select("id,title,view").Order("view desc").Limit(10).Scan(&al).Error;

	if  err != nil {
		return al, err
	}

	if err := CacheSetArticleList(al); err != nil {
		return al , err
	}

	return al, nil
}


// ------------------------ Cache init ------------------------

func CacheInit()  {

}

// TODO 同步cache 和 db
func CacheSync()  {

}

// ------------------------ Cache ------------------------

// 更新 文章访问db
func (a Article) ViewArticle() (err error) {
	if err = db.Model(&a).Update("view", gorm.Expr("view + ?", 1)).Error; err != nil {
		return err
	}
	return err
}

// 设置文章缓存
func (a Article) CacheSetArticle()  (err error) {
	articleKey := e.CACHE_ARTICLE + strconv.Itoa(a.ID)

	if err := gredis.Set(articleKey, a, config.AppSetting.CacheDuration); err != nil {
		return err
	}

	return
}

// 查询文章缓存
func CacheArticle(articleId int) (a Article, err error) {
	articleKey := e.CACHE_ARTICLE + strconv.Itoa(articleId)
	adata, err := gredis.Get(articleKey)

	// redis 查询到key
	if err == nil {
		json.Unmarshal(adata, &a)
		return a, nil
	}
	return a, err
}

// 更新缓存文章访问量
func CacheSetArticleView(a Article) error {
	articleKey := e.CACHE_ARTICLE + strconv.Itoa(a.ID)
	adata, err := gredis.Get(articleKey)

	if err != nil {
	    return err
	}

	json.Unmarshal(adata, &a)
	a.View = a.View + 1
	if err := gredis.Set(articleKey, a, config.AppSetting.CacheDuration); err != nil {
		return err
	}
	return nil
}

// 设置文章列表缓存
func CacheSetArticleList(al []ArticleList) (error) {
	return gredis.Set("current_article_list", al, 3600 * 2)
}

func CacheQueryArticleList() (as []ArticleList, err error) {
	data, err := gredis.Get("current_article_list")
	if err != nil {
	    return as, err
	}

	json.Unmarshal(data, &as)

	return as, nil
}