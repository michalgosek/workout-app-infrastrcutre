package query

import (
	"context"
)

type UserHandlerRepository interface {
	User(ctx context.Context, UUID string) (User, error)
}

type User struct {
	UUID  string
	Role  string
	Name  string
	Email string
}

type UserHandler struct {
	repo UserHandlerRepository
}

func (u *UserHandler) Do(ctx context.Context, UUID string) (User, error) {
	user, err := u.repo.User(ctx, UUID)
	if err != nil {
		return User{}, nil
	}
	return user, nil
}

func NewUserHandlerRepository(r UserHandlerRepository) *UserHandler {
	if r == nil {
		panic("nil user handler repository")
	}
	h := UserHandler{repo: r}
	return &h
}
