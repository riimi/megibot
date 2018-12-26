package model

import (
	"github.com/go-redis/redis"
)

var Redis *redis.Client

func InitRedis(source string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     source,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	Redis = client
	return Redis
}
