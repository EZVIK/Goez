package models

import (
	"Goez/pkg/config"
	"Goez/pkg/utils"
	"fmt"
	"github.com/go-redis/redis/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

var Repo repository

func NewModelTime() gorm.Model {
	t, _ := utils.GetNowTimeCST()
	return gorm.Model{CreatedAt: t, UpdatedAt: t}
}

type repository struct {
	Redis     *redis.Client
	SqlClient *gorm.DB
}

func Setup() {

	mysqlInit(&Repo)
	redisInit(&Repo)
}

func mysqlInit(repo *repository) {

	var err error
	dbConfig := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DatabaseSetting.User,
		config.DatabaseSetting.Password,
		config.DatabaseSetting.Host,
		config.DatabaseSetting.Name)

	gormClient, err := gorm.Open(mysql.Open(dbConfig), &gorm.Config{})
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	repo.SqlClient = gormClient

	repo.SqlClient.AutoMigrate(Article{}, User{}, Record{}, Tag{})

	sqlDB, err := repo.SqlClient.DB()

	// SetConnMaxIdleTime 设置 空闲连接的存活时间
	sqlDB.SetConnMaxIdleTime(5 * time.Second)

	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(3 * time.Second)

	sqlDB.SetMaxIdleConns(3000)
	sqlDB.SetMaxOpenConns(3000)

	fmt.Println("System ... Mysql database initiated.")
}

func redisInit(repo *repository) {

	repo.Redis = redis.NewClient(&redis.Options{
		Addr:     config.RedisSetting.Host,
		Password: config.RedisSetting.Password, // no password set
		DB:       0,                            // use default DB
	})

	str, err := repo.Redis.Ping().Result()

	if err != nil {
		fmt.Println(str)
		panic(err)
	}
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
