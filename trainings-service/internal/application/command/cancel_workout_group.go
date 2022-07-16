package command

import "context"

type CancelTrainerWorkoutGroupService interface {
	CancelTrainerWorkoutGroup(ctx context.Context, groupUUID, trainerUUID string) error
}

type CancelWorkoutGroupHandler struct {
	service CancelTrainerWorkoutGroupService
}

type CancelWorkoutGroup struct {
	GroupUUID   string
	TrainerUUID string
}

func (c *CancelWorkoutGroupHandler) Do(ctx context.Context, cmd CancelWorkoutGroup) error {
	err := c.service.CancelTrainerWorkoutGroup(ctx, cmd.GroupUUID, cmd.TrainerUUID)
	if err != nil {
		return err
	}
	return nil
}

func NewCancelTrainerWorkoutGroupHandler(s CancelTrainerWorkoutGroupService) *CancelWorkoutGroupHandler {
	if s == nil {
		panic("nil cancel trainer workout group service")
	}
	h := CancelWorkoutGroupHandler{service: s}
	return &h
}
