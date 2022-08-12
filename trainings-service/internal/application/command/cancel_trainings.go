package command

import "context"

type CancelTrainingGroupsRepository interface {
	DeleteTrainingGroups(ctx context.Context, trainerUUID string) error
}

type TrainingGroupParticipantNotification struct {
	TrainingUUID string
	TrainerUUID  string
}

type CancelTrainingGroupsHandler struct {
	repo    CancelTrainingGroupsRepository
	service NotificationService
}

func (c *CancelTrainingGroupsHandler) Do(ctx context.Context, trainerUUID string) error {
	err := c.repo.DeleteTrainingGroups(ctx, trainerUUID)
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
