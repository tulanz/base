package cache

import "context"

type Cache interface {
	Init(opts ...Option) error
	Options() Options
	Get(ctx context.Context, key string, resultPtr interface{}, opts ...ReadOption) bool
	Set(ctx context.Context, key string, value interface{}, opts ...WriteOption) error
	Delete(ctx context.Context, key string, opts ...DeleteOption) error
	Exists(ctx context.Context, key string) bool
}
