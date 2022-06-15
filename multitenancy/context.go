package multitenancy

import (
	"context"

	"github.com/asim/go-micro/v3/metadata"
)

const (
	TenantId = "tenant-id"
)

func FromContext(ctx context.Context) (string, bool) {
	tenantId, ok := ctx.Value(TenantId).(string)
	if ok {
		return tenantId, true
	}

	meta, ok := metadata.FromContext(ctx)
	//zap.L().Info("multitenancy.FromContext", zap.Any("metadata", meta), zap.Bool("ok", ok))
	if ok {
		tenantId, ok := meta.Get(TenantId)
		if ok {
			return tenantId, ok
		}
	}
	return "", false
}

func WithContext(ctx context.Context, tenantId string) context.Context {

	meta, ok := metadata.FromContext(ctx)
	if !ok {
		meta = make(metadata.Metadata)
	}
	meta.Set(TenantId, tenantId)
	//zap.L().Info("multitenancy.WithContext", zap.Any("metadata", meta), zap.Bool("ok", ok))
	ctx = metadata.NewContext(ctx, meta)

	return ctx
}
