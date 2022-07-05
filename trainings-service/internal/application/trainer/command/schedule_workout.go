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
	existingGroup, err := s.repository.QueryTrainerWorkoutGroupWithDate(ctx, w.TrainerUUID, w.Date)
	if err != nil {
		const s = "query trainer UUID: %s workout group with date: %v failed, reason: %v"
		logger.Errorf(s, w.TrainerUUID, w.Date, err)
		return ErrRepositoryFailure
	}
	if existingGroup.Date() == w.Date {
		const s = "attempted to create duplicated workout group with date %s failed for trainerUUID: %s"
		logger.Errorf(s, w.Date.String(), w.TrainerUUID)
		return ErrWorkoutGroupDateDuplicated
	}

	group, err := trainer.NewWorkoutGroup(w.TrainerUUID, w.TrainerName, w.GroupName, w.GroupDesc, w.Date)
	if err != nil {
		logger.Errorf("creating new workout group: %+v failed, reason: %v", w, err)
		return err
	}
	err = s.repository.UpsertTrainerWorkoutGroup(ctx, group)
	if err != nil {
		const s = "upsert group UUID: %s for trainer UUID: %s failed, reason: %v"
		logger.Errorf(s, group.UUID(), w.TrainerUUID, err)
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
