package command

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

type UnassignCustomerHandler struct {
	trainerRepository  TrainerRepository
	customerRepository CustomerRepository
}

type UnassignCustomer struct {
	CustomerUUID string
	GroupUUID    string
	TrainerUUID  string
}

func (a *UnassignCustomerHandler) Do(ctx context.Context, c UnassignCustomer) error {
	logger := logrus.WithFields(logrus.Fields{"Trainer-CMD": "UnregisterCustomerHandler"})
	group, err := a.getWorkoutGroup(ctx, c.GroupUUID, c.TrainerUUID)
	if err != nil {
		logger.Errorf("query workout group UUID: %s for trainerUUID: %s failed, reason: %v", c.GroupUUID, c.TrainerUUID, err)
		return err
	}
	err = a.updateCustomerStatus(ctx, c.CustomerUUID, c.GroupUUID)
	if err != nil {
		return err
	}
	err = a.updateTrainerStatus(ctx, c.CustomerUUID, group)
	if err != nil {
		return err
	}
	return nil
}

func (a *UnassignCustomerHandler) updateTrainerStatus(ctx context.Context, customerUUID string, group trainer.WorkoutGroup) error {
	logger := logrus.WithFields(logrus.Fields{"Trainer-CMD": "UnregisterCustomerHandler", "Method": "updateTrainerStatus"})
	group.UnregisterCustomer(customerUUID)
	err := a.trainerRepository.UpsertTrainerWorkoutGroup(ctx, group)
	if err != nil {
		logger.Errorf("upsert workout group UUID: %s workout day failed reason: %v", group.UUID(), err)
		return ErrRepositoryFailure
	}
	return nil
}

func (a *UnassignCustomerHandler) updateCustomerStatus(ctx context.Context, customerUUID, groupUUID string) error {
	logger := logrus.WithFields(logrus.Fields{"Trainer-CMD": "UnregisterCustomerHandler", "Method": "updateCustomerStatus"})
	customerWorkoutDay, err := a.customerRepository.QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID)
	if err != nil {
		logger.Errorf("query customer UUID: %s workout day with groupUUID: %s failed, reason: %v", customerUUID, groupUUID, err)
		return ErrRepositoryFailure
	}
	var empty customer.WorkoutDay
	if customerWorkoutDay == empty {
		logger.Errorf("customer UUID: %s workout day with groupUUID: %s not exist", customerUUID, groupUUID)
		return ErrResourceNotFound
	}
	err = a.customerRepository.DeleteCustomerWorkoutDay(ctx, customerUUID, customerWorkoutDay.UUID())
	if err != nil {
		logger.Errorf("delete customer UUID: %s workout day with UUID: %s failed, reason: %v", customerUUID, customerWorkoutDay.UUID(), err)
		return ErrRepositoryFailure
	}
	return nil
}

func (a *UnassignCustomerHandler) getWorkoutGroup(ctx context.Context, groupUUID, trainerUUID string) (trainer.WorkoutGroup, error) {
	logger := logrus.WithFields(logrus.Fields{"Trainer-CMD": "UnregisterCustomerHandler-updateCustomerStatus", "Method": "getWorkoutGroup"})
	group, err := a.trainerRepository.QueryTrainerWorkoutGroup(ctx, groupUUID)
	if err != nil {
		logger.Errorf("query workout group UUID: %s for trainerUUID: %s failed, reason: %v", groupUUID, trainerUUID, err)
		return trainer.WorkoutGroup{}, ErrRepositoryFailure
	}
	if group.TrainerUUID() != trainerUUID {
		logger.Errorf("workout group UUID: %s does not belong to trainerUUID: %s", group.UUID(), trainerUUID)
		return trainer.WorkoutGroup{}, ErrWorkoutGroupNotOwner
	}
	return group, nil
}

func NewUnassignCustomerHandler(c CustomerRepository, t TrainerRepository) *UnassignCustomerHandler {
	if c == nil {
		panic("nil customer repository")
	}
	if t == nil {
		panic("nil trainer repository")
	}
	return &UnassignCustomerHandler{
		trainerRepository:  t,
		customerRepository: c,
	}
}
