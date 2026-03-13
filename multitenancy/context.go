package multitenancy

import (
	"context"
	// "github.com/asim/go-micro/v3/metadata"
)

const (
	TenantId = "tenant-id"
)

// func FromContext(ctx context.Context) (string, bool) {
// 	tenantId, ok := ctx.Value(TenantId).(string)
// 	if ok {
// 		return tenantId, true
// 	}
// 	meta, ok := FromContext(ctx)
// 	//zap.L().Info("multitenancy.FromContext", zap.Any("metadata", meta), zap.Bool("ok", ok))
// 	if ok {
// 		tenantId, ok := meta.Get(TenantId)
// 		if ok {
// 			return tenantId, ok
// 		}
// 	}
// 	return "", false
// }

// func WithContext(ctx context.Context, tenantId string) context.Context {
// 	meta, ok := ctx.Value(TenantId).(Metadata)
// 	if !ok {
// 		meta = make(Metadata)
// 	}
// 	meta.Set(TenantId, tenantId)
// 	//zap.L().Info("multitenancy.WithContext", zap.Any("metadata", meta), zap.Bool("ok", ok))
// 	ctx = NewContext(ctx, meta)
// 	return ctx
// }

// Metadata is a mapping from metadata keys to values.
type Metadata map[string]interface{}

// type mdKey struct{}

// Len returns the number of items in md.
func (md Metadata) Len() int {
	return len(md)
}

// Copy returns a copy of md.
func (md Metadata) Copy() Metadata {
	return Join(md)
}

// New creates an Metadata from a given key-value map.
func NewMetadata(m map[string]interface{}) Metadata {
	md := Metadata{}
	for k, val := range m {
		md[k] = val
	}
	return md
}

// Join joins any number of mds into a single Metadata.
// The order of values for each key is determined by the order in which
// the mds containing those values are presented to Join.
func Join(mds ...Metadata) Metadata {
	out := Metadata{}
	for _, md := range mds {
		for k, v := range md {
			out[k] = v
		}
	}
	return out
}

// NewContext creates a new context with md attached.
func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, TenantId, md)
}

// FromContext returns the incoming metacode in ctx if it exists.  The
// returned Metadata should not be modified. Writing to it may cause races.
// Modification should be made to copies of the returned Metadata.
func FromContext(ctx context.Context) (md Metadata, ok bool) {
	md, ok = ctx.Value(TenantId).(Metadata)
	return
}
func WithContext(ctx context.Context, tenantId string) context.Context {
	meta, ok := FromContext(ctx)
	if ok {
		nmd := meta.Copy()
		nmd[TenantId] = tenantId
		return NewContext(context.Background(), nmd)
	}
	return context.Background()
}
