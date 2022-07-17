package command

import "context"

type CancelTrainingGroupsService interface {
	CancelTrainingGroups(ctx context.Context, trainerUUID string) error
}

type CancelTrainingGroupsHandler struct {
	service CancelTrainingGroupsService
}

func (c *CancelTrainingGroupsHandler) Do(ctx context.Context, trainerUUID string) error {
	err := c.service.CancelTrainingGroups(ctx, trainerUUID)
	if err != nil {
		return err
	}
	return nil
}

func NewCancelTrainingGroupsHandler(s CancelTrainingGroupsService) *CancelTrainingGroupsHandler {
	if s == nil {
		panic("nil cancel training groups service")
	}
	h := CancelTrainingGroupsHandler{service: s}
	return &h
}
