package mongo

import (
	"context"
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoDB(config *viper.Viper) (*mongo.Database, error) {

	connectionString := config.GetStringSlice("mongo.addr")
	database := config.GetString("mongo.database")

	option := options.Client()
	rb := bson.NewRegistryBuilder()
	rb.RegisterTypeDecoder(reflect.TypeOf(time.Time{}), bsoncodec.NewTimeCodec(bsonoptions.TimeCodec().SetUseLocalTimeZone(true)))

	option.
		SetAuth(options.Credential{
			AuthSource: config.GetString("mongo.AuthSource"),
			Username:   config.GetString("mongo.Username"),
			Password:   config.GetString("mongo.Password"),
		}).
		SetHosts(connectionString).SetConnectTimeout(500 * time.Millisecond).SetHeartbeatInterval(15 * time.Second).
		SetRegistry(rb.Build())

	client, err := mongo.Connect(context.Background(), option)
	if err != nil {
		log.Fatalf("Error opening connection with mongo database: %s", err)
		return nil, err
	}

	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatalf("Error pinging the mongo database: %s", err)
		return nil, err
	}
	return client.Database(database), nil
}

func NewMongoProvider(config *viper.Viper) (*mongo.Database, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	//链接单节点 mongodb
	connectionString := config.GetStringSlice("mongo.addr")
	database := config.GetString("mongo.database")

	if len(connectionString) == 0 {
		return nil, errors.New("addr is empty")
	}

	// Set client options
	option := options.Client().SetHosts(connectionString)
	option.SetConnectTimeout(500 * time.Millisecond)
	option.SetHeartbeatInterval(15 * time.Second)
	option.SetAuth(options.Credential{
		AuthMechanism:           "MONGODB-X509",
		AuthMechanismProperties: map[string]string{"foo": "bar"},
		AuthSource:              "$external",
		Password:                "supersecurepassword",
		Username:                "admin",
	})
	client, err := mongo.Connect(ctx, option)

	if err != nil {
		return nil, err
	}

	// Check the connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client.Database(database), nil
}
