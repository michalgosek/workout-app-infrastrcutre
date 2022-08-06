package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
)

type CancelTrainingGroupRepository interface {
	QueryTrainingGroup(ctx context.Context, trainingUUID string) (trainings.TrainingGroup, error)
	DeleteTrainerGroup(ctx context.Context, trainingUUID, trainerUUID string) error
}

type CancelTrainingGroupHandler struct {
	repo CancelTrainingGroupRepository
}

type CancelWorkoutGroup struct {
	TrainingUUID string
	TrainerUUID  string
}

func (c *CancelTrainingGroupHandler) Do(ctx context.Context, cmd CancelWorkoutGroup) error {
	training, err := c.repo.QueryTrainingGroup(ctx, cmd.TrainingUUID)
	if err != nil {
		return err
	}
	if !training.IsOwnedByTrainer(cmd.TrainerUUID) {
		return ErrTrainingNotOwnedByTrainer
	}

	err = c.repo.DeleteTrainerGroup(ctx, cmd.TrainingUUID, cmd.TrainerUUID)
	if err != nil {
		return err
	}
	return nil
}

func NewCancelTrainingGroupHandler(r CancelTrainingGroupRepository) *CancelTrainingGroupHandler {
	if r == nil {
		panic("nil cancel training group repository")
	}
	h := CancelTrainingGroupHandler{repo: r}
	return &h
}
