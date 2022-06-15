package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/tulanz/base/cache"
	"go.uber.org/zap"
)

type RedisCache struct {
	options cache.Options
	r       *goredis.Client
}

func NewCache(opts ...cache.Option) cache.Cache {
	options := cache.Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return &RedisCache{
		options: options,
	}
}

func (m *RedisCache) Init(opts ...cache.Option) error {
	for _, o := range opts {
		o(&m.options)
	}

	r, err := m.connect()
	if err != nil {
		return err
	}
	m.r = r
	return nil
}

func (m *RedisCache) Options() cache.Options {
	return m.options
}

func (m *RedisCache) connect() (*goredis.Client, error) {
	addrs := m.options.Context.Value(addrsKey{}).([]string)
	if len(addrs) == 0 {
		addrs = []string{":6379"}
	}

	password, _ := m.options.Context.Value(passwordKey{}).(string)
	db, _ := m.options.Context.Value(dbKey{}).(string)

	zap.L().Info("redis", zap.Strings("addrs", addrs), zap.String("passwd", password), zap.String("db", db))
	ctx := context.Background()

	dbInt, _ := strconv.Atoi(db)
	client := goredis.NewClient(&goredis.Options{
		Addr:         strings.Join(addrs, ","),
		Password:     password, // no password set
		DB:           dbInt,    // use default DB
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
}

func (m *RedisCache) prefix(key string) string {
	if m.options.Prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", m.options.Prefix, key)
}

func (m *RedisCache) String() string {
	return "redis"
}

func (m *RedisCache) Exists(ctx context.Context, key string) bool {
	return m.r.Get(ctx, key) == nil
	v, err := m.r.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	return v > 0
}

func (m *RedisCache) Get(ctx context.Context, key string, resultPtr interface{}, opts ...cache.ReadOption) bool {
	readOpts := cache.ReadOptions{}
	for _, o := range opts {
		o(&readOpts)
	}
	key = m.prefix(key)
	data, err := m.r.Get(ctx, key).Bytes()
	if err != nil {
		return false
	}
	err = json.Unmarshal(data, resultPtr)
	if err != nil {
		return false
	}
	return true
}

func (m *RedisCache) Set(ctx context.Context, key string, value interface{}, opts ...cache.WriteOption) error {
	writeOpts := cache.WriteOptions{}
	for _, o := range opts {
		o(&writeOpts)
	}

	key = m.prefix(key)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if writeOpts.Expiry > 0 {
		_, err = m.r.SetEX(ctx, key, data, writeOpts.Expiry).Result()
	} else {
		_, err = m.r.Set(ctx, key, data, -1).Result()
	}
	return err
}

func (m *RedisCache) Close() error {
	return m.r.Close()
}

func (m *RedisCache) Delete(ctx context.Context, key string, opts ...cache.DeleteOption) error {
	deleteOptions := cache.DeleteOptions{}
	for _, o := range opts {
		o(&deleteOptions)
	}

	key = m.prefix(key)

	_, err := m.r.Del(ctx, key).Result()
	return err
}
