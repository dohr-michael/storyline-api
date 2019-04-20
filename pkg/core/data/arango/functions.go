package arango

import (
	"context"
	"github.com/arangodb/go-driver"
	"github.com/dohr-michael/go-libs/errors"
	"github.com/dohr-michael/storyline-api/config"
)

func Database(ctx context.Context, fn func(ctx context.Context, db driver.Database) error) error {
	a, err := config.Config.Arango()
	if err != nil {
		return MapErrors(err)
	}
	client, err := driver.NewClient(*a.ClientConfig)
	if err != nil {
		return MapErrors(err)
	}
	db, err := client.Database(ctx, a.Database)
	if err != nil {
		return MapErrors(err)
	}
	return MapErrors(fn(ctx, db))
}

func Collection(ctx context.Context, name string, fn func(ctx context.Context, col driver.Collection) error) error {
	return Database(ctx, func(ctx context.Context, db driver.Database) error {
		col, err := db.Collection(ctx, name)
		if err != nil {
			return err
		}
		return fn(ctx, col)
	})
}

func InitCollection(ctx context.Context, name string, collectionType driver.CollectionType) error {
	return Database(ctx, func(ctx context.Context, db driver.Database) error {
		_, err := db.Collection(ctx, name)
		if err != nil && driver.IsNotFound(err) {
			_, err := db.CreateCollection(ctx, name, &driver.CreateCollectionOptions{
				Type: collectionType,
			})
			return err
		} else if err != nil {
			return err
		}
		return nil
	})
}

func MapErrors(err error) error {
	switch {
	case driver.IsNotFound(err):
		return errors.NotFoundError
	}
	return err
}
