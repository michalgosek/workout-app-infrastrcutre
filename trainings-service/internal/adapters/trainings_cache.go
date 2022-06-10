package adapters

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain"
)

/*
	Trainer:
	- Can schedule traninings
	- Can update training
	- Can cancel trainings
	- Can set trainigs limit
	- Can remove user from the training

	User (Member):
	- Can assing to training
	- Can cancel training
*/

type TrainingsCacheRepoistory struct{}

func (t *TrainingsCacheRepoistory) GetAllTrainings(ctx context.Context) ([]domain.Training, error) {
	return nil, nil
}

func (t *TrainingsCacheRepoistory) RemoveUserFromTraining(userUUID string, trainingUUID string) error {
	return nil
}

func (t *TrainingsCacheRepoistory) UpdateTrainingWithUUID(ctx context.Context, training domain.Training) error {
	return nil
}

func (t *TrainingsCacheRepoistory) GetTrainingWithUserUUID(ctx context.Context, userUUID string) (domain.Training, error) {
	return domain.Training{}, nil
}

func (t *TrainingsCacheRepoistory) UpdatePlaceLimit(ctx context.Context, n int) error {
	return nil
}

func (t *TrainingsCacheRepoistory) ScheduleTraining(ctx context.Context, training domain.Training) error {
	return nil
}

func (t *TrainingsCacheRepoistory) CancelTraining(ctx context.Context, trainingUUID string) error {
	return nil
}
