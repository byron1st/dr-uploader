package lib

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DefaultDBTransactionTimeout = 5 * time.Second

type DBConfig struct {
	Uri          string
	DatabaseName string
}

var db *mongo.Database

func ConnectDB(uri string, name string) error {
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("%s/%s", uri, name))

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return errors.WithStack(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return errors.WithStack(err)
	}

	db = client.Database(name)
	return nil
}
