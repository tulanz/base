package providers

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedisProvider(vip *viper.Viper) (*redis.Client, error) {

	if !vip.IsSet("cache") {
		return nil, errors.New("")
	}
	config := vip.Sub("cache")
	driver := config.GetString("driver")

	switch driver {
	case "redis":
		ctx := context.Background()
		var client *redis.Client
		client = redis.NewClient(&redis.Options{
			Addr:         config.GetString("addr"),
			Password:     config.GetString("password"), // no password set
			DB:           config.GetInt("db"),          // use default DB
			DialTimeout:  10 * time.Second,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			PoolSize:     10,
			PoolTimeout:  30 * time.Second,
		})
		_, err := client.Ping(ctx).Result()
		if err != nil {
			return nil, err
		}
		return client, nil
	default:
		return nil, errors.New("不支持")
	}
}
