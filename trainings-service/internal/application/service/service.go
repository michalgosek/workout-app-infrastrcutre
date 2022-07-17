package service

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
)

type Command interface {
	InsertTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error
	UpdateTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error
	DeleteTrainingGroup(ctx context.Context, groupUUID, trainerUUID string) error
	DeleteTrainingGroups(ctx context.Context, trainerUUID string) error
}

type Queries interface {
	QueryTrainingGroup(ctx context.Context, trainingUUID, trainerUUID string) (trainings.TrainingGroup, error)
}

type Repository interface {
	Command
	Queries
}

type Trainings struct {
	repo Repository
}

func (t *Trainings) CreateTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error {
	got, err := t.repo.QueryTrainingGroup(ctx, g.UUID(), g.Trainer().UUID())
	if !trainings.IsErrResourceNotFound(err) {
		return err
	}
	if got.IsTrainingDateDuplicated(g.Date()) {
		return errors.New("training group with specified date already exists")
	}

	err = t.repo.InsertTrainingGroup(ctx, g)
	if err != nil {
		return err
	}
	return nil
}

func (t *Trainings) AssignParticipant(ctx context.Context, groupUUID, trainerUUID string, p trainings.Participant) error {
	group, err := t.repo.QueryTrainingGroup(ctx, groupUUID, trainerUUID)
	if err != nil {
		return err
	}
	if group.UUID() == "" {
		return nil
	}
	err = group.AssignParticipant(p)
	if err != nil {
		return err
	}

	err = t.repo.UpdateTrainingGroup(ctx, &group)
	if err != nil {
		return err
	}
	return nil
}

func (t *Trainings) UnassignParticipant(ctx context.Context, trainingUUID, trainerUUID, participantUUID string) error {
	group, err := t.repo.QueryTrainingGroup(ctx, trainingUUID, trainerUUID)
	if err != nil {
		return err
	}
	if group.UUID() == "" {
		return nil
	}
	err = group.UnassignParticipant(participantUUID)
	if err != nil {
		return err
	}

	err = t.repo.UpdateTrainingGroup(ctx, &group)
	if err != nil {
		return err
	}
	return nil
}

func (t *Trainings) CancelTrainingGroup(ctx context.Context, groupUUID, trainerUUID string) error {
	err := t.repo.DeleteTrainingGroup(ctx, groupUUID, trainerUUID)
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
