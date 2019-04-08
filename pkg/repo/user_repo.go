package repo

import (
	"context"
	"github.com/dohr-michael/storyline-api/pkg/model"
)

type UserRepo struct{}

func (r *UserRepo) FindByName(name string, ctx context.Context) (*model.User, error) {
	return nil, nil
}

func (r *UserRepo) FindByEmail(email string, ctx context.Context) (*model.User, error) {
	return nil, nil
}

func (r *UserRepo) CreateUser(name string, email string, ctx context.Context) (*model.User, error) {
	return nil, nil
}
