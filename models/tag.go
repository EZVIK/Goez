package models

import (
	"Goez/pkg/e"
	"errors"
	"time"
)

type Tag struct {
	ID        int        `gorm:"primary_key" json:"id"`
	Name      string     `json:"name"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type ArticleTag struct {
	ID        int `gorm:"primary_key" json:"id"`
	ArticleId int `json:"article_id"`
	TagId     int `json:"tag_id"`
}

// 添加文章绑定标签
func (a *Article) ArtBindTags(ArtId int, tagsId []int) error {

	var at []ArticleTag
	for _, value := range tagsId {
		at = append(at, ArticleTag{ArticleId: ArtId, TagId: value})
	}

	if err := Repo.SqlClient.Create(&at).Error; err != nil {
		return err
	}

	return nil
}

func (t Tag) AddTag() (Tag, error) {
	if t.Name == "" {
		return t, errors.New(e.GetMsg(e.ERROR_TAG_MISSING_PARAMS))
	}
	err := Repo.SqlClient.Create(&t).Error

	if err != nil {
		return t, errors.New(e.GetMsg(e.ERROR_TAG_CREATED_FAILS))
	}
	return t, nil
}

func AddTags(t []Tag) error {
	if err := Repo.SqlClient.Debug().Create(&t).Error; err != nil {
		return err
	}
	return nil
}

func (t Tag) GetTag() (tag Tag, err error) {
	err = Repo.SqlClient.Model(Tag{}).Where("name = ?", t.Name).First(&tag).Error
	return
}

func GetTagsFromRecord(userId int) (tags []Tag, err error) {

	err = Repo.SqlClient.Debug().Model(&Tag{}).Select("tags.id, tags.name").
		Joins("LEFT JOIN article_tags ats ON tags.id = ats.tag_id").
		Joins("LEFT JOIN records re ON ats.article_id = re.article_id").
		Where("user_id = ?", userId).Group("tags.id").Find(&tags).Error

	return
}

func (t Tag) DeleteTag() (err error) {
	err = Repo.SqlClient.Delete(&t).Error

	return
}

func GetTagsFromArticle(article int) (tags []Tag, err error) {

	err = Repo.SqlClient.Debug().Model(&Tag{}).Select("tags.id, tags.name").
		Joins("LEFT JOIN article_tags ats ON tags.id = ats.tag_id").Error
	//Where("ats.article_id = ?", article).Find(&tags).Error

	return
}

func GetTagsByIds(ids []int) (t []Tag, err error) {
	err = Repo.SqlClient.Find(&t, ids).Error
	return
}
