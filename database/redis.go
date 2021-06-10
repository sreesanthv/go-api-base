package database

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Redis struct {
	ctx    context.Context
	rdb    *redis.Client
	logger *logrus.Logger
}

func NewRedis(logger *logrus.Logger) *Redis {
	return &Redis{
		ctx:    context.Background(),
		logger: logger,
		rdb: redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis_addr"),
			Password: viper.GetString("redis_password"), // no password set
			DB:       0,                                 // use default DB
		}),
	}
}

func (r *Redis) Set(key string, value interface{}, expiry time.Duration) {
	err := r.rdb.Set(r.ctx, key, value, expiry).Err()
	if err != nil {
		r.logger.Error("Error writing Redis:", err)
	}
}

func (r *Redis) Get(key string) (string, bool) {
	hasValue := true
	res, err := r.rdb.Get(r.ctx, key).Result()
	if err == redis.Nil {
		hasValue = false
	} else if err != nil {
		hasValue = false
		r.logger.Error("Error accessing Redis:", err)
	}

	return res, hasValue
}
