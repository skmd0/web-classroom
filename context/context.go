package context

import (
	"context"
	"wiki/models/users"
)

const (
	userKey privateKey = "user"
)

type privateKey string

func WithUser(ctx context.Context, user *users.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *users.User {
	if u := ctx.Value(userKey); u != nil {
		if user, ok := u.(*users.User); ok {
			return user
		}
	}
	return nil
}
