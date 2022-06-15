package mongo

import (
	"context"
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

	connectionString := config.GetString("mongo.url")
	database := config.GetString("mongo.database")

	option := options.Client().ApplyURI(connectionString)
	rb := bson.NewRegistryBuilder()
	rb.RegisterTypeDecoder(reflect.TypeOf(time.Time{}), bsoncodec.NewTimeCodec(bsonoptions.TimeCodec().SetUseLocalTimeZone(true)))
	option.SetRegistry(rb.Build())

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
