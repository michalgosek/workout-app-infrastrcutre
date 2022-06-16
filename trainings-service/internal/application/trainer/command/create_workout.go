package command

import (
	"context"
	"fmt"
	"time"

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
	group, err := trainer.NewWorkoutGroup(w.TrainerUUID, w.Name, w.Desc, w.Date)
	if err != nil {
		return "", err
	}
	err = c.repository.UpsertWorkoutGroup(ctx, *group)
	if err != nil {
		const msg = "upsert group UUID: %s for trainer UUID: %s failed, reason: %w"
		return "", fmt.Errorf(msg, group.UUID(), w.TrainerUUID, err)
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
