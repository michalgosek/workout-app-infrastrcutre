package service

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
)

type Command interface {
	InsertTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error
	UpdateTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error
	DeleteTrainingGroup(ctx context.Context, trainingUUID, trainerUUID string) error
	DeleteTrainingGroups(ctx context.Context, trainerUUID string) error
}

type Queries interface {
	QueryTrainingGroup(ctx context.Context, trainingUUID string) (trainings.TrainingGroup, error)
}

type Repository interface {
	Command
	Queries
	IsTrainingGroupDuplicated(ctx context.Context, g *trainings.TrainingGroup) (bool, error)
}

type Trainings struct {
	repo Repository
}

func (t *Trainings) CreateTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error {
	duplicate, err := t.repo.IsTrainingGroupDuplicated(ctx, g)
	if duplicate {
		return ErrTrainingDuplicated
	}
	err = t.repo.InsertTrainingGroup(ctx, g)
	if err != nil {
		return err
	}
	return nil
}

func (t *Trainings) AssignParticipant(ctx context.Context, trainingUUID, trainerUUID string, p trainings.Participant) error {
	training, err := t.repo.QueryTrainingGroup(ctx, trainingUUID)
	if err != nil {
		return err
	}
	if !training.IsOwnedByTrainer(trainerUUID) {
		return ErrTrainingNotOwnedByTrainer
	}
	err = training.AssignParticipant(p)
	if err != nil {
		return err
	}

	err = t.repo.UpdateTrainingGroup(ctx, &training)
	if err != nil {
		return err
	}
	return nil
}

func (t *Trainings) UnassignParticipant(ctx context.Context, trainingUUID, trainerUUID, participantUUID string) error {
	training, err := t.repo.QueryTrainingGroup(ctx, trainingUUID)
	if err != nil {
		return err
	}
	if !training.IsOwnedByTrainer(trainerUUID) {
		return ErrTrainingNotOwnedByTrainer
	}

	err = training.UnassignParticipant(participantUUID)
	if err != nil {
		return err
	}
	err = t.repo.UpdateTrainingGroup(ctx, &training)
	if err != nil {
		return err
	}
	return nil
}

func (t *Trainings) CancelTrainingGroup(ctx context.Context, trainingUUID, trainerUUID string) error {
	training, err := t.repo.QueryTrainingGroup(ctx, trainingUUID)
	if err != nil {
		return err
	}
	if !training.IsOwnedByTrainer(trainerUUID) {
		return ErrTrainingNotOwnedByTrainer
	}

	err = t.repo.DeleteTrainingGroup(ctx, trainingUUID, trainerUUID)
	if err != nil {
		return err
	}
	return nil
}

func (t *Trainings) CancelTrainingGroups(ctx context.Context, trainerUUID string) error {
	err := t.repo.DeleteTrainingGroups(ctx, trainerUUID)
	if err != nil {
		return err
	}
	return nil
}

func NewTrainingsService(r Repository) *Trainings {
	if r == nil {
		panic("nil trainings repository")
	}
	s := Trainings{repo: r}
	return &s
}

var (
	ErrTrainingDuplicated        = errors.New("training group duplicated")
	ErrTrainingNotOwnedByTrainer = errors.New("training not owned by trainer")
)
