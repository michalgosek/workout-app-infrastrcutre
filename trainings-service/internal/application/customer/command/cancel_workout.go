package command

import (
	"context"

	"github.com/sirupsen/logrus"
)

type CancelWorkout struct {
	CustomerUUID string
	TrainerUUUID string
	GroupUUID    string
}

type CancelWorkoutHandler struct {
	trainerRepository  TrainerRepository
	customerRepository CustomerRepository
}

func (c *CancelWorkoutHandler) Do(ctx context.Context, w CancelWorkout) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "CancelWorkoutHandler"})
	group, err := c.trainerRepository.QueryTrainerWorkoutGroup(ctx, w.TrainerUUUID, w.GroupUUID)
	if err != nil {
		logger.Errorf("query workout group UUID: %s failed, reason: %v", w.GroupUUID, err)
		return ErrRepositoryFailure
	}
	//fixme if group not exist then check if it still exists in customer workouts ...

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
