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

type WorkoutUnregister struct {
	CustomerUUID string
	GroupUUID    string
	TrainerUUID  string
}

func (a *UnassignCustomerHandler) Do(ctx context.Context, args WorkoutUnregister) error {
	group, err := a.getWorkoutGroup(ctx, args)
	if err != nil {
		return err
	}
	err = a.updateCustomerStatus(ctx, args)
	if err != nil {
		return err
	}
	err = a.updateTrainerStatus(ctx, args, group)
	if err != nil {
		return err
	}
	return nil
}

func (a *UnassignCustomerHandler) updateTrainerStatus(ctx context.Context, args WorkoutUnregister, group trainer.WorkoutGroup) error {
	logger := logrus.WithFields(logrus.Fields{"Trainer-CMD": "UnregisterCustomerHandler", "Method": "updateTrainerStatus"})

	group.UnregisterCustomer(args.CustomerUUID)
	err := a.trainerRepository.UpsertTrainerWorkoutGroup(ctx, group)
	if err != nil {
		logger.Errorf("upsert workout group UUID: %s workout day failed reason: %v", group.UUID(), err)
		return ErrRepositoryFailure
	}
	return nil
}

func (a *UnassignCustomerHandler) updateCustomerStatus(ctx context.Context, args WorkoutUnregister) error {
	logger := logrus.WithFields(logrus.Fields{"Trainer-CMD": "UnregisterCustomerHandler", "Method": "updateCustomerStatus"})
	customerWorkoutDay, err := a.customerRepository.QueryCustomerWorkoutDay(ctx, args.CustomerUUID, args.GroupUUID)
	if err != nil {
		logger.Errorf("query customer UUID: %s workout day with groupUUID: %s failed, reason: %v", args.CustomerUUID, args.GroupUUID, err)
		return ErrRepositoryFailure
	}
	var empty customer.WorkoutDay
	if customerWorkoutDay == empty {
		logger.Errorf("customer UUID: %s workout day with groupUUID: %s not exist", args.CustomerUUID, args.GroupUUID)
		return ErrResourceNotFound
	}
	err = a.customerRepository.DeleteCustomerWorkoutDay(ctx, args.CustomerUUID, customerWorkoutDay.UUID())
	if err != nil {
		logger.Errorf("delete customer UUID: %s workout day with UUID: %s failed, reason: %v", args.CustomerUUID, customerWorkoutDay.UUID(), err)
		return ErrRepositoryFailure
	}
	return nil
}

func (a *UnassignCustomerHandler) getWorkoutGroup(ctx context.Context, args WorkoutUnregister) (trainer.WorkoutGroup, error) {
	logger := logrus.WithFields(logrus.Fields{"Trainer-CMD": "UnregisterCustomerHandler-updateCustomerStatus", "Method": "getWorkoutGroup"})
	group, err := a.trainerRepository.QueryTrainerWorkoutGroup(ctx, args.GroupUUID)
	if err != nil {
		logger.Errorf("query workout group UUID: %s for trainerUUID: %s failed, reason: %v", args.GroupUUID, args.TrainerUUID, err)
		return trainer.WorkoutGroup{}, ErrRepositoryFailure
	}
	if group.TrainerUUID() != args.TrainerUUID {
		logger.Errorf("workout group UUID: %s does not belong to trainerUUID: %s", group.UUID(), args.TrainerUUID)
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
