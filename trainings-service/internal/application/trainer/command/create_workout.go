package command

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type WorkoutUpserter interface {
	UpsertWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error
}

type CreateWorkoutHandler struct {
	repository WorkoutUpserter
}

type WorkoutGroup struct {
	TrainerUUID string
	Name        string
	Desc        string
	Date        time.Time
}

func (c *CreateWorkoutHandler) Do(ctx context.Context, w WorkoutGroup) (string, error) {
	logger := logrus.WithFields(logrus.Fields{"Component": "CreateWorkoutHandler"})
	group, err := trainer.NewWorkoutGroup(w.TrainerUUID, w.Name, w.Desc, w.Date)
	if err != nil {
		const s = "creating new workout group: %+v failed, reason: %v"
		logger.Errorf(s, w, err)
		return "", err
	}
	err = c.repository.UpsertWorkoutGroup(ctx, *group)
	if err != nil {
		const s = "upsert group UUID: %s for trainer UUID: %s failed, reason: %v"
		logger.Errorf(s, group.UUID(), w.TrainerUUID, err)
		return "", ErrRepositoryFailure
	}
	return group.UUID(), nil
}

func NewCreateWorkoutHandler(r WorkoutUpserter) *CreateWorkoutHandler {
	if r == nil {
		panic("nil repository")
	}
	return &CreateWorkoutHandler{
		repository: r,
	}
}
