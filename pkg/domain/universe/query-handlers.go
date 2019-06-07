package universe

import (
	"context"
	"fmt"
	"github.com/dohr-michael/go-libs/errors"
	"github.com/dohr-michael/go-libs/filters"
	"github.com/dohr-michael/storyline-api/pkg/core"
	"github.com/dohr-michael/storyline-api/pkg/core/data"
	"github.com/dohr-michael/storyline-api/pkg/core/data/arango"
	"github.com/dohr-michael/storyline-api/pkg/domain/relation"
)

type QueryHandlers interface {
	FetchOne(id string, ctx context.Context) (*Universe, error)
	FetchMany(query *filters.Query, ctx context.Context) (*PagedUniverse, error)
	FetchTags(ctx context.Context) ([]Tag, error)
}

type arangoQueryHandlers struct{}

// Search commands
const fetchTagsQuery = `
LET a = (FOR c IN @@collection RETURN c.tags[*])
FOR t IN a[**]
    SORT t ASC
    RETURN DISTINCT { name: t }
`

func (h *arangoQueryHandlers) FetchTags(ctx context.Context) ([]Tag, error) {
	// Check authorizations
	userContext := core.GetUserContext(ctx)
	if userContext == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	d, err := arango.RunQuery(
		ctx,
		arango.NewRunQuery(fetchTagsQuery).WithParam("@collection", Collection),
	)
	if err != nil {
		return nil, err
	}
	var res []Tag
	if err := data.Decode(d, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (h *arangoQueryHandlers) FetchOne(id string, ctx context.Context) (*Universe, error) {
	// Check authorizations
	userContext := core.GetUserContext(ctx)
	if userContext == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	filter, err := filters.NewQuery(
		filters.WithFilter(data.EqStr("_key", id)),
		filters.WithLimit(1),
	)
	if err != nil {
		return nil, err
	}

	paged, err := h.FetchMany(filter, ctx)
	if err != nil {
		return nil, err
	}
	universes := paged.Items
	if len(*universes) == 0 {
		return nil, errors.NotFoundError
	}
	return (*universes)[0], nil
}

func (h *arangoQueryHandlers) FetchMany(query *filters.Query, ctx context.Context) (*PagedUniverse, error) {
	userContext := core.GetUserContext(ctx)
	if userContext == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	items := Universes{}
	total, err := arango.Find(ctx,
		&items,
		arango.
			NewFind(Collection).
			WithRelationRequests(
				arango.NewRelationRequest(
					relation.Collection,
					"createdBy",
					arango.OutDirection,
					string(relation.CreatedBy),
				).WithFieldMapping("email", "_key"),
				arango.NewRelationRequest(
					relation.Collection,
					"owner",
					arango.InDirection,
					string(relation.IsOwnerOf),
				).WithFieldMapping("email", "_key"),
			).
			WithSortByDesc("createdAt").
			WithQuery(query),
	)
	if err != nil {
		return nil, err
	}
	return &PagedUniverse{
		Items:  &items,
		Total:  total,
		Limit:  query.Pager.Limit,
		Offset: query.Pager.Offset,
	}, nil
}
