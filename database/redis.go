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

func (r *Redis) Set(key string, value interface{}, expiry time.Duration) error {
	err := r.rdb.Set(r.ctx, key, value, expiry).Err()
	if err != nil {
		r.logger.Error("Error writing Redis:", err)
	}

	return err
}

func (r *Redis) Get(key string) (string, error) {
	res, err := r.rdb.Get(r.ctx, key).Result()
	if err != nil && err != redis.Nil {
		r.logger.Error("Error accessing Redis:", err)
	}

	return res, err
}

func (r *Redis) Delete(key string) error {
	err := r.rdb.Del(r.ctx, key).Err()
	if err != nil && err != redis.Nil {
		r.logger.Error("Error deleting Redis:", err)
	}
	return nil
}
