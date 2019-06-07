package entity

import (
	"github.com/arangodb/go-driver"
	"github.com/dohr-michael/storyline-api/pkg/core/data/arango"
)

const Collection = "Entities"

type Handlers interface {
	MutationHandlers
	QueryHandlers
}

func NewHandlers() (Handlers, error) {
	err := arango.InitCollection(nil, Collection, driver.CollectionTypeDocument)
	if err != nil {
		return nil, err
	}
	return &arangoHandlers{
		arangoQueryHandlers:    &arangoQueryHandlers{},
		arangoMutationHandlers: &arangoMutationHandlers{},
	}, nil
}

type arangoHandlers struct {
	*arangoQueryHandlers
	*arangoMutationHandlers
}
