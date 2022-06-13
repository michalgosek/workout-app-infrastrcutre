package application

import (
	"context"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain"
	"github.com/sirupsen/logrus"
)

type TrainerRepository interface {
	UpsertSchedule(ctx context.Context, schedule domain.TrainerSchedule) error
	QuerySchedule(ctx context.Context, scheduleUUID string) (domain.TrainerSchedule, error)
	QuerySchedules(ctx context.Context, trainerUUID string) ([]domain.TrainerSchedule, error)
	CancelSchedules(ctx context.Context, scheduleUUIDs ...string) ([]domain.TrainerSchedule, error)
	CancelSchedule(ctx context.Context, scheduleUUID string) (domain.TrainerSchedule, error)
}

type TrainerService struct {
	repository TrainerRepository
}

type TrainerScheduleArgs struct {
	TrainerUUID string
	Name        string
	Desc        string
	Date        time.Time
}

func (t *TrainerService) CreateTrainerSchedule(ctx context.Context, args TrainerScheduleArgs) error {
	schedule, err := domain.NewTrainerSchedule(args.TrainerUUID, args.Name, args.Desc, args.Date)
	if err != nil {
		return fmt.Errorf("creating trainer schedule failed: %v", err)
	}
	err = t.repository.UpsertSchedule(ctx, *schedule)
	if err != nil {
		return fmt.Errorf("upsert schedule failed: %v", err)
	}
	logrus.WithFields(logrus.Fields{"Component": "TrainerService", "Method": "CreateTrainerSchedule"}).Info()
	return nil
}

func (t *TrainerService) GetSchedule(ctx context.Context, trainerUUID string) (domain.TrainerSchedule, error) {
	return domain.TrainerSchedule{}, nil
}

func (t *TrainerService) GetSchedules(ctx context.Context, trainerUUID string) ([]domain.TrainerSchedule, error) {
	return nil, nil
}

func (t *TrainerService) DeleteSchedule(ctx context.Context, sessionUUID, trainerUUID string) error {
	return nil
}

func (t *TrainerService) DeleteSchedules(ctx context.Context, sessionUUID string) error {
	return nil
}

func NewTrainerService(repository TrainerRepository) *TrainerService {
	return &TrainerService{
		repository: repository,
	}
}
