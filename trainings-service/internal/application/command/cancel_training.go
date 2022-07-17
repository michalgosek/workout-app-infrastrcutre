package command

import "context"

type CancelTrainingGroupService interface {
	CancelTrainingGroup(ctx context.Context, trainingUUID, trainerUUID string) error
}

type CancelTrainingGroupHandler struct {
	service CancelTrainingGroupService
}

type CancelWorkoutGroup struct {
	TrainingUUID string
	TrainerUUID  string
}

func (c *CancelTrainingGroupHandler) Do(ctx context.Context, cmd CancelWorkoutGroup) error {
	err := c.service.CancelTrainingGroup(ctx, cmd.TrainingUUID, cmd.TrainerUUID)
	if err != nil {
		return err
	}
	return nil
}

func NewCancelTrainingGroupHandler(s CancelTrainingGroupService) *CancelTrainingGroupHandler {
	if s == nil {
		panic("nil cancel training group service")
	}
	h := CancelTrainingGroupHandler{service: s}
	return &h
}
