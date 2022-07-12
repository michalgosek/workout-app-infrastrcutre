package command

import (
	"context"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainer"
	"github.com/sirupsen/logrus"
	"time"
)

type ScheduleWorkoutArgs struct {
	TrainerUUID string
	TrainerName string
	GroupName   string
	GroupDesc   string
	Date        time.Time
}

type ScheduleWorkoutHandler struct {
	trainerService TrainerService
}

func (s *ScheduleWorkoutHandler) Do(ctx context.Context, args ScheduleWorkoutArgs) error {
	logger := logrus.WithFields(logrus.Fields{"Trainer-CMD": "CreateWorkoutHandler"})
	err := s.trainerService.CreateWorkoutGroup(ctx, trainer.CreateWorkoutGroupArgs{
		TrainerUUID: args.TrainerUUID,
		TrainerName: args.TrainerName,
		GroupName:   args.GroupName,
		GroupDesc:   args.GroupDesc,
		Date:        args.Date,
	})
	if err != nil {
		logger.Errorf("CMD - create trainer workout group failure: %s", err)
		return fmt.Errorf("trainer service failure: %w", err)
	}
	return nil
}

func NewScheduleWorkoutHandler(t TrainerService) (*ScheduleWorkoutHandler, error) {
	if t == nil {
		return nil, ErrNilTrainerService
	}
	h := &ScheduleWorkoutHandler{
		trainerService: t,
	}
	return h, nil
}
