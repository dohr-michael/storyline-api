package user

import (
	"context"
	"fmt"
	"github.com/dohr-michael/storyline-api/pkg/core/data"
	"github.com/dohr-michael/storyline-api/pkg/core/data/arango"
)

type MutationHandlers interface {
	Save(payload *Save, ctx context.Context) (*User, error)
}

type arangoMutationHandlers struct{}

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

func (h *arangoMutationHandlers) Save(payload *Save, ctx context.Context) (*User, error) {
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
