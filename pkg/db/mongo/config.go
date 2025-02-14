package mongo

import (
	"context"
	"os"

	"github.com/go-pkgz/lgr"
)

type MongoConfig struct {
	Mongo struct {
		URI      string `long:"url" env:"URL" default:"http://127.0.0.1:27017" description:"monogo connection url"`
		Database string `long:"database" env:"DATABASE" default:"Application" description:"monogo database name"`
	} `group:"mongo" namespace:"mongo" env-namespace:"MONGO"`
}

func (c MongoConfig) ConnectMongoDB(ctx context.Context, log lgr.L) (*Mongo, error) {
	return New(ctx, log, c.Mongo.URI, c.Mongo.Database)
}

func (c MongoConfig) MustConnectMongoDB(ctx context.Context, log lgr.L) *Mongo {
	mongo, err := c.ConnectMongoDB(ctx, log)
	if err != nil {
		lgr.Default().Logf("[ERROR] mongo connection error: %v", err)
		os.Exit(2)
	}

	return mongo
}
