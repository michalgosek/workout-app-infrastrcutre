package command

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/domain"
)

type RegisterHandlerRepository interface {
	InsertUser(ctx context.Context, u *domain.User) error
	QueryUserWithEmail(ctx context.Context, email string) (domain.User, error)
}

type RegisterUser struct {
	UUID  string
	Role  string
	Name  string
	Email string
}

type RegisterUserHandler struct {
	repo RegisterHandlerRepository
}

func (r *RegisterUserHandler) Do(ctx context.Context, cmd RegisterUser) error {
	got, err := r.repo.QueryUserWithEmail(ctx, cmd.Email)
	if err != nil {
		return nil
	}
	if got.IsEmailAlreadyRegistered(cmd.Email) {
		return errors.New("email already registered")
	}

	u, err := domain.NewUser(cmd.UUID, cmd.Role, cmd.Name, cmd.Email)
	if err != nil {
		return err
	}
	err = r.repo.InsertUser(ctx, u)
	if err != nil {
		return nil
	}
	return nil
}

func NewRegisterHandlerRepository(r RegisterHandlerRepository) *RegisterUserHandler {
	if r == nil {
		panic("nil register handler repository")
	}
	h := RegisterUserHandler{repo: r}
	return &h
}
