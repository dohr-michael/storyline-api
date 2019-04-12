package arango

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/dohr-michael/go-libs/errors"
	"github.com/dohr-michael/go-libs/filters"
	"github.com/dohr-michael/go-libs/storage"
	"github.com/dohr-michael/storyline-api/pkg/core/data"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"time"
)

const fetchByIdQuery = `
FOR c IN @@collection
	FILTER c._key == @myId
	LIMIT 0, 1
	RETURN c
`

const fetchManyQuery = `
let total = COUNT(@@collection)

let items = (
	FOR c IN @@collection
        LIMIT @offset, @count
	    RETURN c
)
RETURN {total: total, items: items}
`

const fetchManyWithFilterQuery = `
let total = LENGTH(
	FOR c IN @@collection
		FILTER %s
		RETURN 1
)

let items = (
	FOR c IN @@collection
		FILTER %s
		LIMIT @offset, @count
		RETURN c
)
RETURN {total: total, items: items}
`

var _ = data.Repository(&Repository{})

type LogLevel string

const (
	Debug = LogLevel("DEBUG")
	Info  = LogLevel("INFO")
)

type Repository struct {
	Collection  string
	LogLevel    LogLevel
	OneFactory  func() interface{}
	ManyFactory func() interface{}
}

type RepositoryConfig func(c *Repository) *Repository

func WithLogLevel(level LogLevel) func(c *Repository) *Repository {
	return func(c *Repository) *Repository {
		c.LogLevel = level
		return c
	}
}

func NewRepository(
	collection string,
	collectionType driver.CollectionType,
	oneFactory func() interface{},
	manyFactory func() interface{},
	configs ...RepositoryConfig,
) (*Repository, error) {
	// Initialize collection.
	if err := InitCollection(nil, collection, collectionType); err != nil {
		return nil, err
	}
	res := &Repository{
		Collection:  collection,
		OneFactory:  oneFactory,
		ManyFactory: manyFactory,
		LogLevel:    Info,
	}
	for _, fn := range configs {
		res = fn(res)
	}

	return res, nil
}

func ToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
	}
}

func (r *Repository) Decode(input interface{}, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			ToTimeHookFunc()),
		Result: result,
	})
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}
	return err
}

func (r *Repository) FetchOne(id string, ctx context.Context) (interface{}, error) {
	res := r.OneFactory()
	items, err := RunQuery(
		ctx,
		fetchByIdQuery,
		map[string]interface{}{
			"@collection": r.Collection,
			"myId":        id,
		},
	)

	if err != nil {
		return nil, err
	} else if len(items) == 0 {
		return nil, errors.NotFoundError
	}
	err = r.Decode(items[0], res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repository) FetchMany(query *filters.Query, ctx context.Context) (res *storage.Paged, err error) {
	result := r.ManyFactory()
	res = &storage.Paged{
		Items: result,
		Query: query,
	}
	filterStr, args := RsqlToFilter("c", query.Filter)
	args["@collection"] = r.Collection
	args["offset"] = query.Pager.Offset
	args["count"] = query.Pager.Limit

	queryStr := fetchManyQuery
	if len(filterStr) > 0 {
		queryStr = fmt.Sprintf(fetchManyWithFilterQuery, filterStr, filterStr)
	}
	items, err := RunQuery(
		ctx,
		queryStr,
		args,
	)
	if err != nil {
		return nil, err
	} else if len(items) == 0 {
		return nil, errors.NotFoundError
	}

	doc := items[0]
	s, ok := doc["items"].([]interface{})
	if !ok {
		return nil, errors.NotFoundError
	}
	err = r.Decode(&s, &result)

	if err != nil {
		return nil, err
	}

	res.Total = int64(doc["total"].(float64))

	if err != nil {
		return
	}
	return
}

func (r *Repository) Create(entity interface{}, ctx context.Context) (id string, res interface{}, err error) {
	err = Collection(ctx, r.Collection, func(ctx context.Context, col driver.Collection) error {
		meta, err := col.CreateDocument(ctx, entity)
		if err != nil {
			return err
		}
		id = meta.Key
		return nil
	})
	if err != nil {
		return
	}
	res, err = r.FetchOne(id, ctx)
	return
}

func (r *Repository) Update(id string, toUpdate interface{}, ctx context.Context) (res interface{}, err error) {
	err = Collection(ctx, r.Collection, func(ctx context.Context, col driver.Collection) error {
		_, err := col.UpdateDocument(ctx, id, toUpdate)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	res, err = r.FetchOne(id, ctx)
	return
}

func (r *Repository) Save(id string, entity interface{}, ctx context.Context) (res interface{}, err error) {
	err = Collection(ctx, r.Collection, func(ctx context.Context, col driver.Collection) error {
		exists, err := col.DocumentExists(ctx, id)
		if err != nil {
			return err
		}
		if exists {
			res, err = r.Update(id, entity, ctx)
			return err
		} else {
			_, res, err = r.Create(entity, ctx)
			return err
		}
	})
	return
}

func (r *Repository) Delete(id string, ctx context.Context) (err error) {
	err = Collection(ctx, r.Collection, func(ctx context.Context, col driver.Collection) error {
		_, err := col.RemoveDocument(ctx, id)
		return err
	})
	return
}
