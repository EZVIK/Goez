package models

import (
	"Goez/pkg/config"
	"Goez/pkg/utils"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type ID struct {
	ID int `gorm:"primary_key" json:"id"`
}

func NewModelTime() gorm.Model {
	t, _ := utils.GetNowTimeCST()
	return gorm.Model{CreatedAt: t, UpdatedAt: t}
}

var db *gorm.DB

func Setup() {

	var err error
	dbConfig := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DatabaseSetting.User,
		config.DatabaseSetting.Password,
		config.DatabaseSetting.Host,
		config.DatabaseSetting.Name)

	db, err = gorm.Open(mysql.Open(dbConfig), &gorm.Config{})
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	db.AutoMigrate(Article{}, User{}, Record{}, Tag{})

	// 数据表单数 user users
	//db.SingularTable(true)

	//
	//gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	//	return config.DatabaseSetting.TablePrefix + defaultTableName
	//}

	//db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	//db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	//db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	//db.DB().SetMaxIdleConns(100)
	//db.DB().SetMaxOpenConns(1000)

	sqlDB, err := db.DB()

	// SetConnMaxIdleTime 设置 空闲连接的存活时间
	sqlDB.SetConnMaxIdleTime(5 * time.Second)

	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(3 * time.Second)

	// SetMaxOpenConns 设置打开数据库连接的最大数量
	//sqlDB.SetMaxOpenConns(1000)
	//
	//sqlDB.SetMaxIdleConns(1000)

	fmt.Println("System ... Mysql database initiated.")

}

func PageChecker(pageIndex string, pageSize string) (pi int, ps int) {

	pi, err1 := strconv.Atoi(pageIndex)
	ps, err2 := strconv.Atoi(pageSize)

	if err1 != nil {
		pi = 1
	}

	if err2 != nil {
		ps = 10
	}

	return pi, ps
}

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
