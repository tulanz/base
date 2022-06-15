package redis

import (
	"context"

	"github.com/tulanz/base/cache"
)

type addrsKey struct{}
type passwordKey struct{}
type dbKey struct{}

func WithAddrs(addrs ...string) cache.Option {
	return func(o *cache.Options) {
		o.Context = context.WithValue(o.Context, addrsKey{}, addrs)
	}
}

func WithPassword(password string) cache.Option {
	return func(o *cache.Options) {
		o.Context = context.WithValue(o.Context, passwordKey{}, password)
	}
}

func WithDB(db string) cache.Option {
	return func(o *cache.Options) {
		o.Context = context.WithValue(o.Context, dbKey{}, db)
	}
}
