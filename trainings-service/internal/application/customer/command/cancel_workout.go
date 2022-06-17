package command

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

type CancelWorkoutHandlerRepository interface {
	QueryWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error)
	UpsertWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error
	DeleteCustomerWorkoutDay(ctx context.Context, customerUUID, groupUUID string) error
}

type CancelWorkoutHandler struct {
	repository CancelWorkoutHandlerRepository
}

type CancelWorkoutDetails struct {
	CustomerUUID string
	GroupUUID    string
}

func (c *CancelWorkoutHandler) Do(ctx context.Context, w CancelWorkoutDetails) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "CancelWorkoutHandler"})
	group, err := c.repository.QueryWorkoutGroup(ctx, w.GroupUUID)
	if err != nil {
		logger.Errorf("query workout group UUID: %s failed, reason: %v", w.GroupUUID, err)
		return ErrRepositoryFailure
	}
	err = c.repository.DeleteCustomerWorkoutDay(ctx, w.CustomerUUID, w.GroupUUID)
	if err != nil {
		logger.Errorf("delete customer UUID: %s workout day UUID: %s failed, reason: %v", w.CustomerUUID, w.GroupUUID, err)
		return ErrRepositoryFailure
	}

	group.UnregisterCustomer(w.CustomerUUID)
	err = c.repository.UpsertWorkoutGroup(ctx, group)
	if err != nil {
		logger.Errorf("upsert workout group UUID: %s failed, reason: %v", w.GroupUUID, err)
		return ErrRepositoryFailure
	}
	return nil
}

func NewCancelWorkoutHandler(r CancelWorkoutHandlerRepository) *CancelWorkoutHandler {
	return &CancelWorkoutHandler{repository: r}
}
