package service

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
)

type Command interface {
	InsertTrainerWorkoutGroup(ctx context.Context, g *trainings.WorkoutGroup) error
	UpdateTrainerWorkoutGroup(ctx context.Context, g *trainings.WorkoutGroup) error
	DeleteTrainerWorkoutGroup(ctx context.Context, groupUUID, trainerUUID string) error
}

type Queries interface {
	QueryTrainerWorkoutGroup(ctx context.Context, groupUUID, trainerUUID string) (trainings.WorkoutGroup, error)
}

type Repository interface {
	Command
	Queries
}

type Trainings struct {
	repo Repository
}

func (t *Trainings) CreateTrainerWorkoutGroup(ctx context.Context, g *trainings.WorkoutGroup) error {
	got, err := t.repo.QueryTrainerWorkoutGroup(ctx, g.UUID(), g.Trainer().UUID())
	if err != nil {
		return err
	}
	if got.Date() != g.Date() {
		return errors.New("resource already exist")
	}

	err = t.repo.InsertTrainerWorkoutGroup(ctx, g)
	if err != nil {
		return err
	}
	return nil
}

func (t *Trainings) AssignParticipant(ctx context.Context, groupUUID, trainerUUID string, p trainings.Participant) error {
	group, err := t.repo.QueryTrainerWorkoutGroup(ctx, groupUUID, trainerUUID)
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

	err = t.repo.UpdateTrainerWorkoutGroup(ctx, &group)
	if err != nil {
		return err
	}
	return nil
}

func (t *Trainings) UnassignParticipant(ctx context.Context, groupUUID, trainerUUID, participantUUID string) error {
	group, err := t.repo.QueryTrainerWorkoutGroup(ctx, groupUUID, trainerUUID)
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

	err = t.repo.UpdateTrainerWorkoutGroup(ctx, &group)
	if err != nil {
		return err
	}
	return nil
}

func (t *Trainings) CancelTrainerWorkoutGroup(ctx context.Context, groupUUID, trainerUUID string) error {
	err := t.repo.DeleteTrainerWorkoutGroup(ctx, groupUUID, trainerUUID)
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
