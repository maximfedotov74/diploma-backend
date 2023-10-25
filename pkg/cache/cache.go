package cache

import (
	"context"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var redisOnce sync.Once

func NewCacheClient(addr string, password string) *redis.Client {
	redisOnce.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       0,
		})
		_, err := client.Ping(context.Background()).Result()
		if err != nil {
			log.Fatalf("Cannot connect to redis: %s", err.Error())
		}
		redisClient = client
	})
	return redisClient
}
