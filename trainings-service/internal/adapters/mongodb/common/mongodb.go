package common

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoClient(addr string, timeout time.Duration) (*mongo.Client, error) {
	opts := options.Client()
	opts.ApplyURI(addr)
	opts.SetConnectTimeout(timeout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	mongoCLI, err := mongo.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("mongo client creation failed: %v", err)
	}
	err = mongoCLI.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("mongo client connection failed: %v", err)
	}
	err = mongoCLI.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("mongo client ping req failed: %v", err)
	}
	return mongoCLI, nil
}
