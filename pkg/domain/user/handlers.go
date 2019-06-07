package user

import (
	"github.com/arangodb/go-driver"
	"github.com/dohr-michael/storyline-api/pkg/core/data/arango"
)

const Collection = "Users"

type Handlers interface {
	QueryHandlers
	MutationHandlers
}

func NewHandlers() (Handlers, error) {
	err := arango.InitCollection(nil, Collection, driver.CollectionTypeDocument)
	if err != nil {
		return nil, err
	}
	return &arangoHandler{
		arangoQueryHandlers:    &arangoQueryHandlers{},
		arangoMutationHandlers: &arangoMutationHandlers{},
	}, nil
}

type arangoHandler struct {
	*arangoQueryHandlers
	*arangoMutationHandlers
}
