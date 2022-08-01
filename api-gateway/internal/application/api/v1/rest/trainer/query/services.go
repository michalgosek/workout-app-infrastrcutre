package query

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/v1/users/query"
)

type UsersService interface {
	User(ctx context.Context, UUID string) (query.User, error)
}

type TrainingsService interface {
	TrainingGroup(ctx context.Context, q TrainingQuery) (TrainingGroup, error)
}
