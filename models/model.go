package models

import (
	"Goez/pkg/config"
	"Goez/pkg/utils"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

type ID struct {
	ID         int `gorm:"primary_key" json:"id"`
}

type ModelTime struct {
	CreatedOn  time.Time `gorm:"" json:"created_on"`
	ModifiedOn time.Time `json:"modified_on"`
	DeletedOn  *time.Time `json:"deleted_on"`
}

func NewModelTime() ModelTime {
	t,_ := utils.GetNowTimeCST()

	return ModelTime{CreatedOn: t, ModifiedOn: t}
}

var db *gorm.DB

func Setup()  {

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
	sqlDB.SetMaxOpenConns(0)

	sqlDB.SetMaxIdleConns(0)



	fmt.Println("System ... Mysql database initiated.")

}

// CloseDB closes database connection (unnecessary)
func CloseDB() {

}
