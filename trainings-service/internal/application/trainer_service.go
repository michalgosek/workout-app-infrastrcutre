package application

import (
	"context"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"time"
)

type TrainerRepository interface {
	UpsertWorkoutGroup(ctx context.Context, schedule trainer.WorkoutGroup) error
	QueryWorkoutGroup(ctx context.Context, UUID, trainerUUID string) (trainer.WorkoutGroup, error)
	QueryWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error)
	DeleteWorkoutGroups(ctx context.Context, trainerUUID string) error
	DeleteWorkoutGroup(ctx context.Context, UUID, trainerUUID string) error
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

func (t *TrainerService) CreateWorkoutGroup(ctx context.Context, args TrainerSchedule) (string, error) {
	schedule, err := trainer.NewWorkoutGroup(args.TrainerUUID, args.Name, args.Desc, args.Date)
	if err != nil {
		return "", fmt.Errorf("creating trainer workout group failed: %v", err)
	}
	err = t.repository.UpsertWorkoutGroup(ctx, *schedule)
	if err != nil {
		return "", fmt.Errorf("upsert schedule UUID: %s for trainer UUID: %s failed, reason: %w", schedule.UUID(), args.TrainerUUID, err)
	}
	return schedule.UUID(), nil
}

func (t *TrainerService) GetWorkoutGroup(ctx context.Context, scheduleUUID, trainerUUID string) (trainer.WorkoutGroup, error) {
	schedule, err := t.repository.QueryWorkoutGroup(ctx, scheduleUUID, trainerUUID)
	if err != nil {
		return trainer.WorkoutGroup{}, fmt.Errorf("query workout group failed: %v", err)
	}
	if schedule.TrainerUUID() != trainerUUID {
		return trainer.WorkoutGroup{}, nil
	}
	return schedule, nil
}

func (t *TrainerService) AssignCustomer(ctx context.Context, customerUUID, workoutGroupUUID, trainerUUID string) error {
	schedule, err := t.repository.QueryWorkoutGroup(ctx, workoutGroupUUID, trainerUUID)
	if err != nil {
		return err
	}
	if schedule.TrainerUUID() != trainerUUID {
		return ErrScheduleNotOwner
	}
	err = schedule.AssignCustomer(customerUUID)
	if err != nil {
		return fmt.Errorf("assign customer to the schedule failed: %w", err)
	}
	err = t.repository.UpsertWorkoutGroup(ctx, schedule)
	if err != nil {
		return fmt.Errorf("upsert schedule failed: %w", ErrRepositoryFailure)
	}
	return nil
}

func (t *TrainerService) GetWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error) {
	schedules, err := t.repository.QueryWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		return nil, fmt.Errorf("get schedules failed: %v", err)
	}
	return schedules, nil
}

func (t *TrainerService) DeleteWorkoutGroup(ctx context.Context, scheduleUUID, trainerUUID string) error {
	schedule, err := t.repository.QueryWorkoutGroup(ctx, scheduleUUID, trainerUUID)
	if err != nil {
		return fmt.Errorf("query trainer workout group failed: %v", err)
	}
	if schedule.TrainerUUID() != trainerUUID {
		return ErrScheduleNotOwner
	}
	err = t.repository.DeleteWorkoutGroup(ctx, scheduleUUID, trainerUUID)
	if err != nil {
		return fmt.Errorf("delete workout group failed: %v", err)
	}
	return nil
}

func (t *TrainerService) DeleteWorkoutGroups(ctx context.Context, trainerUUID string) error {
	if trainerUUID == "" {
		return nil
	}
	err := t.repository.DeleteWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		return fmt.Errorf("cancel schedules failed: %v", err)
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
