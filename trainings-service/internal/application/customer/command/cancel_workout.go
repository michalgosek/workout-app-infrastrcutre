package command

import (
	"context"

	"github.com/sirupsen/logrus"
)

type CancelWorkoutDetails struct {
	CustomerUUID string
	GroupUUID    string
}

type CancelWorkoutHandler struct {
	trainerRepository  TrainerRepository
	customerRepository CustomerRepository
}

func (c *CancelWorkoutHandler) Do(ctx context.Context, w CancelWorkoutDetails) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "CancelWorkoutHandler"})
	group, err := c.trainerRepository.QueryTrainerWorkoutGroup(ctx, w.GroupUUID)
	if err != nil {
		logger.Errorf("query workout group UUID: %s failed, reason: %v", w.GroupUUID, err)
		return ErrRepositoryFailure
	}
	err = c.customerRepository.DeleteCustomerWorkoutDay(ctx, w.CustomerUUID, w.GroupUUID)
	if err != nil {
		logger.Errorf("delete customer UUID: %s workout day UUID: %s failed, reason: %v", w.CustomerUUID, w.GroupUUID, err)
		return ErrRepositoryFailure
	}

	group.UnregisterCustomer(w.CustomerUUID)
	err = c.trainerRepository.UpsertTrainerWorkoutGroup(ctx, group)
	if err != nil {
		logger.Errorf("upsert workout group UUID: %s failed, reason: %v", w.GroupUUID, err)
		return ErrRepositoryFailure
	}
	return nil
}

func NewCancelWorkoutHandler(c CustomerRepository, t TrainerRepository) *CancelWorkoutHandler {
	if c == nil {
		panic("nil customer repository")
	}
	if t == nil {
		panic("nil trainer repository")
	}
	return &CancelWorkoutHandler{trainerRepository: t, customerRepository: c}
}
