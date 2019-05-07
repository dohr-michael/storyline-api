package universe

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/dohr-michael/go-libs/errors"
	"github.com/dohr-michael/go-libs/filters"
	"github.com/dohr-michael/go-libs/storage"
	"github.com/dohr-michael/storyline-api/pkg/core"
	"github.com/dohr-michael/storyline-api/pkg/core/data"
	"github.com/dohr-michael/storyline-api/pkg/core/data/arango"
	"github.com/dohr-michael/storyline-api/pkg/domain/relation"
	"github.com/dohr-michael/storyline-api/pkg/domain/user"
	"github.com/google/uuid"
)

const Collection = "Universes"

type Handlers interface {
	storage.ReadRepository
	FetchTags(ctx context.Context) ([]Tag, error)
	// Write
	Create(payload *Create, ctx context.Context) (*Universe, error)
}

func NewHandlers() (Handlers, error) {
	err := arango.InitCollection(nil, Collection, driver.CollectionTypeDocument)
	if err != nil {
		return nil, err
	}
	return &arangoHandler{}, nil
}

type arangoHandler struct{}

// Search commands
const fetchTagsQuery = `
LET a = (FOR c IN @@collection RETURN c.tags[*])
FOR t IN a[**]
    SORT t ASC
    RETURN DISTINCT { name: t }
`

func (h *arangoHandler) FetchTags(ctx context.Context) ([]Tag, error) {
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

func (h *arangoHandler) FetchOne(id string, ctx context.Context) (storage.Entity, error) {
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
	universes := paged.Items.(*Universes)
	if len(*universes) == 0 {
		return nil, errors.NotFoundError
	}
	return (*universes)[0], nil
}

func (h *arangoHandler) FetchMany(query *filters.Query, ctx context.Context) (*storage.Paged, error) {
	userContext := core.GetUserContext(ctx)
	if userContext == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	return arango.Find(ctx,
		&Universes{},
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
}

// Mutation commands
const createQuery = `
FOR u IN @@collection_users
    FILTER u._key == @creator
    INSERT MERGE({createdAt:  DATE_NOW()}, @data) INTO @@collection
    LET inserted = NEW
    LET relations = [
        { _from: u._id, _to: inserted._id, kind: @is_owner_of },
        { _from: inserted._id, _to: u._id, kind: @created_by }
    ]
    LET ignored = (FOR r IN relations INSERT r INTO @@collection_relations)
    RETURN MERGE(inserted, {createdBy: u, owner: u})
`

func (h *arangoHandler) Create(payload *Create, ctx context.Context) (*Universe, error) {
	u := core.GetUserContext(ctx)
	if u == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	d, err := data.ToMap(payload)
	if err != nil {
		return nil, err
	}
	d["_key"] = uuid.New().String()
	results, err := arango.RunQuery(ctx,
		arango.NewRunQuery(createQuery).
			WithParam("@collection_users", user.Collection).
			WithParam("@collection_relations", relation.Collection).
			WithParam("@collection", Collection).
			WithParam("creator", u.Email).
			WithParam("created_by", relation.CreatedBy).
			WithParam("is_owner_of", relation.IsOwnerOf).
			WithParam("data", d))
	if err != nil {
		return nil, err
	} else if len(results) == 0 {
		// User not found
		return nil, fmt.Errorf("unknow connected user")
	}
	res := &Universe{}
	if err := data.Decode(results[0], res); err != nil {
		return nil, err
	}
	return res, nil
}
