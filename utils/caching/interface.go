package caching

import "time"

type CacheRepository interface {
	Get(key string) ([]byte, error)
	Set(key string, entry []byte, expiration time.Duration) error
}
