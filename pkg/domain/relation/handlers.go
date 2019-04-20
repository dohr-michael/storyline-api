package relation

import (
	"github.com/arangodb/go-driver"
	"github.com/dohr-michael/storyline-api/pkg/core/data/arango"
)

const Collection = "Relations"

type Handlers interface{}

func NewHandlers() (Handlers, error) {
	err := arango.InitCollection(nil, Collection, driver.CollectionTypeEdge)
	if err != nil {
		return nil, err
	}
	return &arangoHandlers{}, nil
}

type arangoHandlers struct{}
