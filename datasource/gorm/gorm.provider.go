package gorm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/tulanz/base/datasource"
	"github.com/tulanz/base/multitenancy"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGormTenancy(config *viper.Viper, entityMap datasource.EntityMap) (multitenancy.Tenancy, error) {
	driver := config.GetString("db.driver")
	connectionString := config.GetString("db.connection_string")
	isolation := config.Get("multitenancy.isolation")
	if isolation == "" {
		isolation = "schema"
	}
	if len(driver) == 0 {
		return nil, errors.New("driver is empty")
	}

	if len(connectionString) == 0 {
		return nil, errors.New("connection_string is empty")
	}

	defaultDB, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	defaultDB.Debug()

	var clientCreateFn func(ctx context.Context, tenantId string) (multitenancy.Resource, error)
	var clientCloseFunc func(resource multitenancy.Resource)
	switch isolation {
	case "schema":
		clientCreateFn = func(ctx context.Context, tenantId string) (multitenancy.Resource, error) {
			db := defaultDB // gorm.Open(driver, connectionString)
			//autoMigrate(tenantId, entityMap, db)
			return db, nil
		}
		clientCloseFunc = func(resource multitenancy.Resource) {}
	case "database":
		dbName := "" // 从连接中获取
		clientCreateFn = func(ctx context.Context, tenantId string) (multitenancy.Resource, error) {
			dsn := strings.Replace(connectionString, dbName, DBName(dbName, tenantId), 1)
			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // gorm.Open(driver, dsn)
			if err != nil {
				return nil, err
			}
			db.Debug()
			return db, nil
		}
		clientCloseFunc = func(resource multitenancy.Resource) {}
	}
	tenancy := multitenancy.NewCachedTenancy(clientCreateFn, clientCloseFunc)
	return tenancy, nil
}

// DBName returns the prefixed database name in order to avoid collision with MySQL internal databases.
func DBName(prefix string, tenantId string) string {
	if len(tenantId) == 0 {
		return prefix
	}
	return fmt.Sprintf("%s_%s", prefix, tenantId)
}

// FromDBName returns the source name of the tenant.
func FromDBName(serviceName string, name string) string {
	return strings.TrimPrefix(name, fmt.Sprintf("%s_", serviceName))
}

func DBFromContext(tenancy multitenancy.Tenancy, ctx context.Context) (*gorm.DB, error) {
	tenantName, _ := multitenancy.FromContext(ctx)
	db, err := tenancy.ResourceFor(ctx, tenantName)
	if err != nil {
		return nil, err
	}
	return db.(*gorm.DB), nil
}

func TableName(tableName string, tenantId string) string {
	if len(tenantId) == 0 {
		return tableName
	}
	return fmt.Sprintf("%s_%s", tableName, tenantId)
}

// func autoMigrate(tenantId string, entityMap datasource.EntityMap, db *gorm.DB) error {
// 	entities := entityMap.GetEntities()
// 	db = db.Unscoped()
// 	for _, entity := range entities {
// 		scope := db.NewScope(entity)
// 		tableName := TableName(scope.TableName(), tenantId)
// 		db = db.Table(tableName).AutoMigrate(entity)
// 	}
// 	return nil
// }
