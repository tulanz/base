package elasticsearch

import (
	"context"

	"github.com/tulanz/elastic/v7"
	"github.com/tulanz/base/datasource"
	"github.com/tulanz/base/multitenancy"
	"go.uber.org/zap"
)

type tenancy struct {
}

func NewElasticSearchTenancy(client *elastic.Client, logger *zap.Logger, entityMap datasource.EntityMap) multitenancy.Tenancy {
	var clientCreateFn = func(ctx context.Context, tenantId string) (multitenancy.Resource, error) {
		err := autoMigrate(ctx, client, logger, entityMap, tenantId)
		if err != nil {
			return nil, err
		}
		return client, nil
	}

	var clientCloseFunc = func(resource multitenancy.Resource) {

	}
	return multitenancy.NewCachedTenancy(clientCreateFn, clientCloseFunc)
}

func autoMigrate(c context.Context, client *elastic.Client, logger *zap.Logger, entityMap datasource.EntityMap, tenantId string) error {
	for _, entity := range entityMap.GetEntities() {
		indexModel := NewIndexModel(client, logger, entity, tenantId)
		if err := indexModel.CreateIndex(c); err != nil {
			return err
		}
	}
	return nil
}
