package redis

import (
	"context"
	"github.com/capstone-kelompok-7/backend-disappear/config"
	"github.com/capstone-kelompok-7/backend-disappear/utils/caching"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

type redisCacheRepository struct {
	rdb *redis.Client
}

func NewRedisClient(cnf config.Config) caching.CacheRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cnf.Redis.Addr,
		Password: cnf.Redis.Pass,
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logrus.Fatalf("Error koneksi Redis: %s", err.Error())
	}

	logrus.Info("Redis connected successfully")

	return &redisCacheRepository{
		rdb: rdb,
	}
}

func (r redisCacheRepository) Get(key string) ([]byte, error) {
	val, err := r.rdb.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

func (r redisCacheRepository) Set(key string, entry []byte, expiration time.Duration) error {
	return r.rdb.Set(context.Background(), key, entry, expiration).Err()
}
