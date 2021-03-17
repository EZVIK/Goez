package models

import (
	"Goez/pkg/e"
	"github.com/jinzhu/gorm"
)

type Record struct {
	ID        int `gorm:"primary_key" json:"id"`
	UserId    int `gorm:"index" json:"user_id"`
	ArticleId int `gorm:"index" json:"article_id"`
	gorm.Model
}

func AddRecord(r Record) (int, error) {

	// 查询是否已记录过
	if err := db.Where("user_id = ? and article_id = ? ", r.UserId, r.ArticleId).First(&r).Error; err != nil {
		// 除了无法找到记录, 其余异常返回
		if err.Error() != e.ErrorMsg["RECORD_NOT_FOUND"] {
			return 0, err
		}
	} else {
		// 已存在浏览记录则返回 TODO 是否允许多次保存d∂浏览记录
		if r.ID != 0 {
			return -1, nil
		}
	}

	time := NewModelTime()
	r.CreatedAt = time.CreatedAt
	r.UpdatedAt = time.UpdatedAt

	if err := db.Create(&r).Error; err != nil {
		return 0, err
	} else {
		return r.ID, nil
	}
}
