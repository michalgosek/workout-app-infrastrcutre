package query

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

type WorkoutGetter interface {
	QueryWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error)
}

type GetWorkoutHandler struct {
	repository WorkoutGetter
}

func (t *GetWorkoutHandler) Do(ctx context.Context, groupUUID, trainerUUID string) (trainer.WorkoutGroup, error) {
	logger := logrus.WithFields(logrus.Fields{"Component": "GetWorkoutHandler"})

	group, err := t.repository.QueryWorkoutGroup(ctx, groupUUID)
	if err != nil {
		const s = "query workout groupUUID: %s for trainerUUID: %s failed, reason: %v"
		logger.Errorf(s, groupUUID, trainerUUID, err)
		return trainer.WorkoutGroup{}, ErrRepositoryFailure
	}
	if group.UUID() == "" {
		return trainer.WorkoutGroup{}, nil
	}
	if group.TrainerUUID() != trainerUUID {
		const s = "query workout group UUID: %s does not belong to trainerUUID: %s"
		logger.Errorf(s, groupUUID, trainerUUID)
		return trainer.WorkoutGroup{}, ErrWorkoutGroupNotOwner
	}
	return group, nil
}

func NewGetWorkoutHandler(w WorkoutGetter) *GetWorkoutHandler {
	if w == nil {
		panic("nil repository")
	}
	return &GetWorkoutHandler{
		repository: w,
	}
}