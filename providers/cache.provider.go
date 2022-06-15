package providers

import (
	"github.com/spf13/viper"
	"github.com/tulanz/base/cache"
	"github.com/tulanz/base/cache/memory"
	"github.com/tulanz/base/cache/redis"
	"go.uber.org/zap"
)

func NewCacheProvider(config *viper.Viper) cache.Cache {
	driver := config.GetString("cache.driver")
	prefix := config.GetString("cache.prefix")

	var c cache.Cache
	switch driver {
	case "redis":
		addrs := config.GetString("cache.addr")
		passwd := config.GetString("cache.password")
		db := config.GetString("cache.db")
		c = redis.NewCache(cache.WithPrefix(prefix), redis.WithAddrs(addrs), redis.WithPassword(passwd), redis.WithDB(db))
	default:
		c = memory.NewCache(cache.WithPrefix(prefix))
	}

	return c
}

func InitCache(c cache.Cache, logger *zap.Logger) {
	if err := c.Init(); err != nil {
		logger.Error("cache init", zap.Error(err))
	}
}
