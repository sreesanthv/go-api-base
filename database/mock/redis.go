package mock

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisMock struct {
	store map[string]string
}

func NewRedis() *RedisMock {
	store := make(map[string]string)
	return &RedisMock{store}
}

func (r *RedisMock) Set(key string, value string, expiry time.Duration) error {
	r.store[key] = value
	return nil
}

func (r *RedisMock) Get(key string) (string, error) {
	var err error
	val, ok := r.store[key]
	if !ok {
		err = redis.Nil
	}
	return val, err
}

func (r *RedisMock) Delete(key string) error {
	delete(r.store, key)
	return nil
}
