package mongo

import (
	"context"

	"github.com/go-pkgz/lgr"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(ctx context.Context, log lgr.L, uri, db string) (*Mongo, error) {
	cl, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &Mongo{cl: cl, db: db, log: log}, nil
}

type Mongo struct {
	log lgr.L

	cl *mongo.Client
	db string
}

func (m *Mongo) Operation(fn func(db *mongo.Database) error) error {
	db := m.cl.Database(m.db)
	return fn(db)
}

func (m *Mongo) Session(ctx context.Context, actions ...func(sessionCtx mongo.SessionContext, db *mongo.Database) error) error {
	session, err := m.cl.StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sessionCtx mongo.SessionContext) error {
		for _, action := range actions {
			err := action(sessionCtx, m.cl.Database(m.db))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (m *Mongo) Transaction(ctx context.Context, actions ...func(ctx context.Context, db *mongo.Database) error) error {
	return m.Session(ctx, func(sessionCtx mongo.SessionContext, db *mongo.Database) error {
		err := sessionCtx.StartTransaction()
		if err != nil {
			return err
		}

		for i, action := range actions {
			err = action(ctx, db)
			if err != nil {
				m.log.Logf("[ERROR] transaction action %v error: %v", i, err)
				break
			}
		}

		if err != nil {
			err = sessionCtx.AbortTransaction(ctx)
			if err != nil {
				m.log.Logf("[ERROR] abort transaction error: %v", err)
				return err
			}
		}

		err = sessionCtx.CommitTransaction(ctx)
		if err != nil {
			m.log.Logf("[ERROR] commit transaction error: %v", err)
			return err
		}

		sessionCtx.EndSession(ctx)
		return nil
	})
}

func (m *Mongo) Disconnect(ctx context.Context) {
	err := m.cl.Disconnect(ctx)
	if err != nil {
		m.log.Logf("[ERROR] mongo disconnect error: %v", err)
	}
}
