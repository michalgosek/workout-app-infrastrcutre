package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"time"
)

type ScheduleTrainerWorkoutGroup struct {
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	TrainerUUID string    `json:"trainer_uuid"`
	TrainerName string    `json:"trainer_name"`
}

type CreateTrainerWorkoutGroupService interface {
	CreateTrainerWorkoutGroup(ctx context.Context, g *trainings.WorkoutGroup) error
}

type ScheduleTrainerWorkoutGroupHandler struct {
	service CreateTrainerWorkoutGroupService
}

func (s *ScheduleTrainerWorkoutGroupHandler) Do(ctx context.Context, cmd ScheduleTrainerWorkoutGroup) error {
	t, err := trainings.NewTrainer(cmd.TrainerUUID, cmd.TrainerName)
	if err != nil {
		return err
	}
	g, err := trainings.NewWorkoutGroup(cmd.UUID, cmd.Name, cmd.Description, cmd.Date, t)
	if err != nil {
		return err
	}

	err = s.service.CreateTrainerWorkoutGroup(ctx, g)
	if err != nil {
		return err
	}
	return nil
}

func NewScheduleTrainerWorkoutGroupHandler(s CreateTrainerWorkoutGroupService) *ScheduleTrainerWorkoutGroupHandler {
	if s == nil {
		panic("nil create training workout group service")
	}
	h := ScheduleTrainerWorkoutGroupHandler{service: s}
	return &h
}
