package universe

import (
	"context"
	"fmt"
	"github.com/dohr-michael/storyline-api/pkg/core"
	"github.com/dohr-michael/storyline-api/pkg/core/data"
	"github.com/dohr-michael/storyline-api/pkg/core/data/arango"
	"github.com/dohr-michael/storyline-api/pkg/domain/relation"
	"github.com/dohr-michael/storyline-api/pkg/domain/user"
	"github.com/google/uuid"
)

type MutationHandlers interface {
	Create(payload *Create, ctx context.Context) (*Universe, error)
}

type arangoMutationHandlers struct {}

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

func (h *arangoMutationHandlers) Create(payload *Create, ctx context.Context) (*Universe, error) {
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