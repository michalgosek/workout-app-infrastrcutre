package adapters

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/domain"
)

type UsersCacheRepoistory struct{}

// User:
// - Can sing up for trainings
// - Can cancel his trainings
// - Can see all trainings

func (t *UsersCacheRepoistory) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return nil, nil
}

func (t *UsersCacheRepoistory) FindUserWithUUID(ctx context.Context, UUID string) ([]domain.User, error) {
	return nil, nil
}
