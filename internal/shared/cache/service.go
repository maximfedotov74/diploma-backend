package cache

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	redis *redis.Client
	ctx   context.Context
}

var redisClient *redis.Client
var redisOnce sync.Once

func NewCacheService(addr string, password string, ctx context.Context) *CacheService {
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
	return &CacheService{
		redis: redisClient,
		ctx:   ctx,
	}
}

func (cs CacheService) Set(key string, value interface{}, exp time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return cs.redis.Set(cs.ctx, key, p, exp).Err()
}

func (cs *CacheService) Get(key string, dest interface{}) error {
	result, err := cs.redis.Get(cs.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(result), dest)
}

func (cs *CacheService) Shutdown() error {
	return cs.redis.ShutdownSave(cs.ctx).Err()
}
