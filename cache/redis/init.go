package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisOpts struct {
	Endpoint string
	Username string
	Password string
	Database int
}

var redisClient *redis.Client

func Redis() *redis.Client {
	if redisClient == nil {
		panic("redis 模块没有初始化")
	}

	return redisClient
}

func NewRedisClient(redisOpts *RedisOpts) error {
	client := redis.NewClient(&redis.Options{
		Addr:     redisOpts.Endpoint,
		Username: redisOpts.Username,
		Password: redisOpts.Password,
		Protocol: 3,
		DB:       redisOpts.Database,
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return err
	}
	redisClient = client
	return nil
}
