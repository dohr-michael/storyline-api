package graphql

import "context"

func (*Resolver) UserByName(ctx context.Context, args *struct{ Name string }) (*userResolver, error) {
	return nil, nil
}

func (*Resolver) UserByEmail(ctx context.Context, args *struct{ Email string }) (*userResolver, error) {
	return nil, nil
}
