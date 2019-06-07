package graphql

import (
	"context"
	errors2 "github.com/dohr-michael/go-libs/errors"
	"github.com/dohr-michael/go-libs/filters"
	"github.com/dohr-michael/storyline-api/pkg/core/data"
	"github.com/graphql-go/graphql"
)

func byIdQuery(
	fieldId string,
	baseType *graphql.Object,
	fn func(string, context.Context) (interface{}, error),
) *graphql.Field {
	return &graphql.Field{
		Type: baseType,
		Args: graphql.FieldConfigArgument{
			fieldId: &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			res, err := fn(p.Args[fieldId].(string), p.Context)
			if err == errors2.NotFoundError {
				return nil, nil
			} else if err != nil {
				return nil, err
			}
			return res, nil
		},
	}
}

func pagedQuery(
	prefix string,
	baseType *graphql.Object,
	fn func(*filters.Query, context.Context) (interface{}, error),
) *graphql.Field {
	return &graphql.Field{
		Type: pagedType(prefix, baseType),
		Args: graphql.FieldConfigArgument{
			"query": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"limit": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"offset": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			var conf []filters.QueryConfig
			if v, ok := p.Args["query"]; ok {
				conf = append(conf, filters.WithRsqlFilter(v.(string)))
			}
			if v, ok := p.Args["limit"]; ok {
				conf = append(conf, filters.WithLimit(int64(v.(int))))
			}
			if v, ok := p.Args["offset"]; ok {
				conf = append(conf, filters.WithOffset(int64(v.(int))))
			}

			query, err := filters.NewQuery(conf...)
			if err != nil {
				return nil, err
			}
			return fn(query, p.Context)
		},
	}
}

func pagedType(prefix string, baseType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        prefix + "Paged",
		Description: "",
		Fields: graphql.Fields{
			"items": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(baseType))),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if paged, ok := p.Source.(data.Paged); ok {
						return paged.GetItems(), nil
					}
					return nil, nil
				},
			},
			"total": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if paged, ok := p.Source.(data.Paged); ok {
						return paged.GetTotal(), nil
					}
					return nil, nil
				},
			},
			"limit": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if paged, ok := p.Source.(data.Paged); ok {
						return paged.GetLimit(), nil
					}
					return nil, nil
				},
			},
			"offset": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if paged, ok := p.Source.(data.Paged); ok {
						return paged.GetOffset(), nil
					}
					return nil, nil
				},
			},
		},
	})
}
