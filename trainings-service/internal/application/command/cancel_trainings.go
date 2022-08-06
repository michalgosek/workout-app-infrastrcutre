package command

import "context"

type CancelTrainingGroupsRepository interface {
	DeleteTrainerGroups(ctx context.Context, trainerUUID string) error
}

type CancelTrainingGroupsHandler struct {
	repo CancelTrainingGroupsRepository
}

func (c *CancelTrainingGroupsHandler) Do(ctx context.Context, trainerUUID string) error {
	err := c.repo.DeleteTrainerGroups(ctx, trainerUUID)
	if err != nil {
		return err
	}
	return nil
}

func NewCancelTrainingGroupsHandler(r CancelTrainingGroupsRepository) *CancelTrainingGroupsHandler {
	if r == nil {
		panic("nil cancel training groups repository")
	}
	h := CancelTrainingGroupsHandler{repo: r}
	return &h
}
