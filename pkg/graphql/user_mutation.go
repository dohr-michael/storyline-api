package graphql

import "context"

func (*Resolver) CreateUser(ctx context.Context, args *struct {
	Name  string
	Email string
}) (*userResolver, error) {
	return nil, nil
}
