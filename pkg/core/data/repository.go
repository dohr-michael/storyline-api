package data

import "github.com/dohr-michael/go-libs/storage"

type Repository interface {
	storage.ReadRepository
	storage.WriteRepository
}
