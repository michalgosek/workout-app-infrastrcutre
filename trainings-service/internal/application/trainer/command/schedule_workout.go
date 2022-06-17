package command

import (
	"context"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

type ScheduleWorkoutHandlerRepository interface {
	UpsertWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error
}

type ScheduleWorkoutHandler struct {
	repository ScheduleWorkoutHandlerRepository
}

type WorkoutGroupDetails struct {
	TrainerUUID string
	Name        string
	Desc        string
	Date        time.Time
}

func (s *ScheduleWorkoutHandler) Do(ctx context.Context, w WorkoutGroupDetails) (string, error) {
	logger := logrus.WithFields(logrus.Fields{"Component": "CreateWorkoutHandler"})
	group, err := trainer.NewWorkoutGroup(w.TrainerUUID, w.Name, w.Desc, w.Date)
	if err != nil {
		const s = "creating new workout group: %+v failed, reason: %v"
		logger.Errorf(s, w, err)
		return "", err
	}

	err = s.repository.UpsertWorkoutGroup(ctx, *group)
	if err != nil {
		const s = "upsert group UUID: %s for trainer UUID: %s failed, reason: %v"
		logger.Errorf(s, group.UUID(), w.TrainerUUID, err)
		return "", ErrRepositoryFailure
	}
	return group.UUID(), nil
}

func NewScheduleWorkoutHandler(r ScheduleWorkoutHandlerRepository) *ScheduleWorkoutHandler {
	if r == nil {
		panic("nil repository")
	}
	return &ScheduleWorkoutHandler{
		repository: r,
	}
}
