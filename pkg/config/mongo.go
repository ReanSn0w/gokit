package config

import (
	"context"
	"os"

	"github.com/ReanSn0w/gokit/pkg/db/mongo"
	"github.com/go-pkgz/lgr"
)

type MongoConfig struct {
	Mongo struct {
		URI      string `long:"url" env:"URL" default:"http://127.0.0.1:27017" description:"monogo connection url"`
		Database string `long:"database" env:"DATABASE" default:"Application" description:"monogo database name"`
	} `group:"mongo" namespace:"mongo" env-namespace:"MONGO"`
}

func (c MongoConfig) MongoConnect(ctx context.Context, log lgr.L) (*mongo.Mongo, error) {
	return mongo.New(ctx, log, c.Mongo.URI, c.Mongo.Database)
}

func (c MongoConfig) MongoMustConnect(ctx context.Context, log lgr.L) *mongo.Mongo {
	mongo, err := c.MongoConnect(ctx, log)
	if err != nil {
		lgr.Default().Logf("[ERROR] mongo connection error: %v", err)
		os.Exit(2)
	}

	return mongo
}
