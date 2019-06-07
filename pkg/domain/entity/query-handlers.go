package entity

import "github.com/dohr-michael/go-libs/filters"

type QueryHandlers interface {
	FindByUniverse(universeId string, query filters.Query)
}

type arangoQueryHandlers struct {}

