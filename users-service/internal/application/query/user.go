package query

import (
	"context"
)

type UserHandlerRepository interface {
	User(ctx context.Context, UUID string) (User, error)
}

type User struct {
	Role  string `json:"role"`
	Name  string `json:"name"`
	Email string `json:"email"`
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
