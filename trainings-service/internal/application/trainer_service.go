package application

import (
	"context"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type TrainerRepository interface {
	UpsertSchedule(ctx context.Context, schedule trainer.TrainerSchedule) error
	QuerySchedule(ctx context.Context, scheduleUUID string) (trainer.TrainerSchedule, error)
	QuerySchedules(ctx context.Context, trainerUUID string) ([]trainer.TrainerSchedule, error)
	CancelSchedules(ctx context.Context, scheduleUUIDs ...string) ([]trainer.TrainerSchedule, error)
	CancelSchedule(ctx context.Context, scheduleUUID string) (trainer.TrainerSchedule, error)
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

func (t *TrainerService) CreateTrainerSchedule(ctx context.Context, args TrainerSchedule) error {
	schedule, err := trainer.NewSchedule(args.TrainerUUID, args.Name, args.Desc, args.Date)
	if err != nil {
		return fmt.Errorf("creating trainer schedule failed: %v", err)
	}
	err = t.repository.UpsertSchedule(ctx, *schedule)
	if err != nil {
		return fmt.Errorf("upsert schedule UUID: %s for trainer UUID: %s failed, reason: %w", schedule.UUID(), args.TrainerUUID, err)
	}
	return nil
}

func (t *TrainerService) GetSchedule(ctx context.Context, scheduleUUID, trainerUUID string) (trainer.TrainerSchedule, error) {
	schedule, err := t.repository.QuerySchedule(ctx, scheduleUUID)
	if err != nil {
		return trainer.TrainerSchedule{}, fmt.Errorf("query schedule failed: %v", err)
	}
	if schedule.TrainerUUID() != trainerUUID {
		return trainer.TrainerSchedule{}, nil
	}
	return schedule, nil
}

func (t *TrainerService) AssingCustomer(ctx context.Context, customerUUID, scheduleUUID, trainerUUID string) error {
	schedule, err := t.repository.QuerySchedule(ctx, scheduleUUID)
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
	err = t.repository.UpsertSchedule(ctx, schedule)
	if err != nil {
		return fmt.Errorf("upsert schedule failed: %w", ErrRepositoryFailure)
	}
	return nil
}

func (t *TrainerService) GetSchedules(ctx context.Context, trainerUUID string) ([]trainer.TrainerSchedule, error) {
	schedules, err := t.repository.QuerySchedules(ctx, trainerUUID)
	if err != nil {
		return nil, fmt.Errorf("get schedules failed: %v", err)
	}
	return schedules, nil
}

func (t *TrainerService) DeleteSchedule(ctx context.Context, sessionUUID, trainerUUID string) error {
	return nil
}

func (t *TrainerService) DeleteSchedules(ctx context.Context, sessionUUID string) error {
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
