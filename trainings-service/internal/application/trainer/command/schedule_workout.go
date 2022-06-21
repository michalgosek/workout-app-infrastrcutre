package command

import (
	"context"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

type ScheduleWorkoutHandler struct {
	repository TrainerRepository
}

type ScheduleWorkout struct {
	TrainerUUID string
	TrainerName string
	GroupName   string
	GroupDesc   string
	Date        time.Time
}

func (s *ScheduleWorkoutHandler) Do(ctx context.Context, w ScheduleWorkout) error {
	logger := logrus.WithFields(logrus.Fields{"Trainer-CMD": "CreateWorkoutHandler"})
	group, err := trainer.NewWorkoutGroup(w.TrainerUUID, w.TrainerName, w.GroupName, w.GroupDesc, w.Date)
	if err != nil {
		logger.Errorf("creating new workout group: %+v failed, reason: %v", w, err)
		return err
	}
	err = s.repository.UpsertTrainerWorkoutGroup(ctx, *group)
	if err != nil {
		logger.Errorf("upsert group UUID: %s for trainer UUID: %s failed, reason: %v", group.UUID(), w.TrainerUUID, err)
		return ErrRepositoryFailure
	}
	return nil
}

func NewScheduleWorkoutHandler(t TrainerRepository) *ScheduleWorkoutHandler {
	if t == nil {
		panic("nil trainer repository")
	}
	return &ScheduleWorkoutHandler{
		repository: t,
	}
}
