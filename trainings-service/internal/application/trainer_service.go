package application

import (
	"context"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type TrainerCommands interface {
	UpsertWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error
	DeleteWorkoutGroups(ctx context.Context, trainerUUID string) error
	DeleteWorkoutGroup(ctx context.Context, groupUUID string) error
}

type TrainerQueries interface {
	QueryWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error)
	QueryWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error)
}

type TrainerRepository interface {
	TrainerCommands
	TrainerQueries
}

type TrainerService struct {
	repository TrainerRepository
}

type TrainerSchedule struct {
	TrainerUUID string
	Name        string
	Desc        string
	Date        time.Time
}

type WorkoutRegistration struct {
	CustomerUUID string
	TrainerUUID  string
	GroupUUID    string
}

func (t *TrainerService) CreateWorkoutGroup(ctx context.Context, args TrainerSchedule) (string, error) {
	group, err := trainer.NewWorkoutGroup(args.TrainerUUID, args.Name, args.Desc, args.Date)
	if err != nil {
		return "", fmt.Errorf("creating workout group failed - invalid arguments data: %v", err)
	}
	err = t.repository.UpsertWorkoutGroup(ctx, *group)
	if err != nil {
		return "", fmt.Errorf("upsert group UUID: %s for trainer UUID: %s failed, reason: %w", group.UUID(), args.TrainerUUID, err)
	}
	return group.UUID(), nil
}

func (t *TrainerService) GetWorkoutGroup(ctx context.Context, groupUUID, trainerUUID string) (trainer.WorkoutGroup, error) {
	group, err := t.repository.QueryWorkoutGroup(ctx, groupUUID)
	if err != nil {
		return trainer.WorkoutGroup{}, fmt.Errorf("query workout group failed: %v", err)
	}
	if group.TrainerUUID() != trainerUUID {
		return trainer.WorkoutGroup{}, nil
	}
	return group, nil
}

func (t *TrainerService) AssignCustomer(ctx context.Context, args WorkoutRegistration) error {
	group, err := t.repository.QueryWorkoutGroup(ctx, args.GroupUUID)
	if err != nil {
		return err
	}
	if group.TrainerUUID() != args.TrainerUUID {
		return ErrScheduleNotOwner
	}
	err = group.AssignCustomer(args.CustomerUUID)
	if err != nil {
		return fmt.Errorf("assign customer to the group failed: %w", err)
	}
	err = t.repository.UpsertWorkoutGroup(ctx, group)
	if err != nil {
		return fmt.Errorf("upsert group failed: %w", ErrRepositoryFailure)
	}
	return nil
}

func (t *TrainerService) UnregisterCustomer(ctx context.Context, args WorkoutRegistration) error {
	return nil
}

func (t *TrainerService) GetWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error) {
	groups, err := t.repository.QueryWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		return nil, fmt.Errorf("get groups failed: %v", err)
	}
	return groups, nil
}

func (t *TrainerService) DeleteWorkoutGroup(ctx context.Context, groupUUID, trainerUUID string) error {
	group, err := t.repository.QueryWorkoutGroup(ctx, groupUUID)
	if err != nil {
		return fmt.Errorf("query trainer workout group failed: %v", err)
	}
	if group.TrainerUUID() != trainerUUID {
		return ErrScheduleNotOwner
	}
	err = t.repository.DeleteWorkoutGroup(ctx, groupUUID)
	if err != nil {
		return fmt.Errorf("delete workout group failed: %v", err)
	}
	return nil
}

func (t *TrainerService) DeleteWorkoutGroups(ctx context.Context, trainerUUID string) error {
	err := t.repository.DeleteWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		return fmt.Errorf("delete workout groups failed: %v", err)
	}
	return nil
}

func NewTrainerService(repository TrainerRepository) *TrainerService {
	if repository == nil {
		panic("repository is nil")
	}
	return &TrainerService{
		repository: repository,
	}
}
