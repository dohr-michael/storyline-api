package user

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
)

const Collection = "Users"

type Handlers interface {
	storage.ReadRepository
	// Write
	Save(payload *Save, ctx context.Context) (*User, error)
}

func NewHandlers() (Handlers, error) {
	err := arango.InitCollection(nil, Collection, driver.CollectionTypeDocument)
	if err != nil {
		return nil, err
	}
	return &arangoHandler{}, nil
}

type arangoHandler struct{}

func (h *arangoHandler) FetchOne(id string, ctx context.Context) (storage.Entity, error) {

	// Check authorizations
	userContext := core.GetUserContext(ctx)
	if userContext == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	filter, err := filters.NewQuery(
		filters.WithFilter(data.EqStr("email", id)),
		filters.WithLimit(1),
	)
	if err != nil {
		return nil, err
	}

	paged, err := h.FetchMany(filter, ctx)
	if err != nil {
		return nil, err
	}
	users := paged.Items.(*Users)
	if len(*users) == 0 {
		return nil, errors.NotFoundError
	}
	return (*users)[0], nil
}

func (h *arangoHandler) FetchMany(query *filters.Query, ctx context.Context) (*storage.Paged, error) {
	// Check authorizations
	userContext := core.GetUserContext(ctx)
	if userContext == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	return arango.Find(ctx,
		&Users{},
		arango.NewFind(Collection).
			WithQuery(query).
			WithFieldId("email"),
	)
}

const saveQuery = `
UPSERT { _key: @key }
INSERT MERGE({ 
	_key: @key, 
	createdAt: DATE_NOW(),
	lastConnection: DATE_NOW()
}, @data)
UPDATE MERGE({ 
	lastConnection: DATE_NOW()
}, @data) IN @@collection
RETURN NEW
`

func (h *arangoHandler) Save(payload *Save, ctx context.Context) (*User, error) {
	if payload.User == nil {
		return nil, fmt.Errorf("unknow connected user")
	}
	results, err := arango.RunQuery(ctx,
		arango.NewRunQuery(saveQuery).
			WithParams(
				map[string]interface{}{
					"@collection": Collection,
					"key":         payload.User.Email,
					"data": map[string]interface{}{
						"name":    payload.User.Name,
						"locale":  payload.User.Locale,
						"gender":  payload.User.Gender,
						"picture": payload.User.Picture,
					},
				},
			),
	)
	if err != nil {
		return nil, err
	} else if len(results) == 0 {
		return nil, nil
	}
	res := &User{}
	err = data.Decode(results[0], res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
