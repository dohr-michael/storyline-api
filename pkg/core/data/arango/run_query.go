package arango

import (
	"context"
	"github.com/arangodb/go-driver"
)

type runQueryParams struct {
	query  string
	params map[string]interface{}
}

func NewRunQuery(query string) *runQueryParams {
	return &runQueryParams{
		query:  query,
		params: map[string]interface{}{},
	}
}

func (q *runQueryParams) WithParams(params map[string]interface{}) *runQueryParams {
	q.params = params
	return q
}

func (q *runQueryParams) WithParam(key string, value interface{}) *runQueryParams {
	q.params[key] = value
	return q
}

func RunQuery(ctx context.Context, params *runQueryParams) (QueryResults, error) {
	results := QueryResults{}
	err := Database(ctx, func(ctx context.Context, db driver.Database) error {
		// log.Printf("Run query %v, %v", query, params)
		cursor, err := db.Query(driver.WithQueryCount(ctx), params.query, params.params)
		if err != nil {
			return err
		}
		for cursor.HasMore() {
			result := make(QueryResult)
			_, err := cursor.ReadDocument(ctx, &result)
			if err != nil {
				return err
			}
			results = append(results, result)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}
