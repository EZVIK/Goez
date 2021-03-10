package models

import (
	"Goez/pkg/e"
	"errors"
)

type Tag struct {
	ID 			int 		`gorm:"primary_key" json:"id"`
	Name    	string		`json:"name"`
}

type ArticleTag struct {
	ID 			int 		`gorm:"primary_key" json:"id"`
	ArticleId	int			`json:"articleId"`
	TagId		int			`json:"tagId"`
}

// 添加文章绑定标签
func (a *Article) ArtBindTags(ArtId int, tagsId []int) error{

	var at []ArticleTag
	for _, value := range tagsId {
		at = append(at, ArticleTag{ArticleId: ArtId, TagId:  value})
	}

	if err := db.Create(&at).Error; err!= nil {
		return err
	}

	return nil
}

func (t Tag) AddTag() (Tag, error) {
	if t.Name == "" {
		return t, errors.New(e.GetMsg(e.ERROR_TAG_MISSING_PARAMS))
	}
	err := db.Create(&t).Error

	if err != nil {
		return t, errors.New(e.GetMsg(e.ERROR_TAG_CREATED_FAILS))
	}
	return t, nil
}

func AddTags(t []Tag) (error) {
	if err := db.Debug().Create(&t).Error; err != nil {
		return err
	}
	return nil
}

func GetTagsById(id int) (t Tag, err error) {
	t.ID = id
	err = db.First(&t).Error
	return
}

func (t Tag) GetTag() (tag Tag, err error) {
	err = db.Model(Tag{}).Where("name = ?", t.Name).First(&tag).Error
	return
}