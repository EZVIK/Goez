package gredis

import (
	"Goez/pkg/config"
	"fmt"
	"github.com/go-redis/redis/v7"
)

var RedisClient *redis.Client

// Setup Initialize the Redis instance
func Setup() *redis.Client {

	fmt.Println("System ...Redis databse config initiated.")

	return redis.NewClient(&redis.Options{
		Addr:     config.RedisSetting.Host,
		Password: config.RedisSetting.Password, // no password set
		DB:       0,                            // use default DB
	})
}
