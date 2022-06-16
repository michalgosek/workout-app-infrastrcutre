package command

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

type CustomerUnregister interface {
	UpsertWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error
	QueryWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error)
	QueryCustomerWorkoutDay(ctx context.Context, customerUUID, trainerWorkoutGroupUUID string) (customer.WorkoutDay, error)
	DeleteCustomerWorkoutDay(ctx context.Context, customerUUID, customerWorkoutDayUUID string) error
}

type UnassignCustomerHandler struct {
	repository CustomerUnregister
}

type WorkoutUnregister struct {
	CustomerUUID string
	GroupUUID    string
	TrainerUUID  string
}

func (a *UnassignCustomerHandler) Do(ctx context.Context, args WorkoutUnregister) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "UnregisterCustomerHandler"})

	group, err := a.repository.QueryWorkoutGroup(ctx, args.GroupUUID)
	if err != nil {
		const s = "query workout group UUID: %s for trainerUUID: %s failed, reason: %v"
		logger.Errorf(s, args.GroupUUID, args.TrainerUUID, err)
		return ErrRepositoryFailure
	}
	if group.TrainerUUID() != args.TrainerUUID {
		const s = "workout group UUID: %s does not belong to trainerUUID: %s"
		logger.Errorf(s, group.UUID(), args.TrainerUUID)
		return ErrWorkoutGroupNotOwner
	}

	customerWorkoutDay, err := a.repository.QueryCustomerWorkoutDay(ctx, args.CustomerUUID, args.GroupUUID)
	if err != nil {
		const s = "query customer UUID: %s workout day with groupUUID: %s failed, reason: %v"
		logger.Errorf(s, args.CustomerUUID, args.GroupUUID, err)
		return ErrRepositoryFailure
	}
	err = a.repository.DeleteCustomerWorkoutDay(ctx, args.CustomerUUID, customerWorkoutDay.UUID())
	if err != nil {
		const s = "delete customer UUID: %s workout day with UUID: %s failed, reason: %v"
		logger.Errorf(s, args.CustomerUUID, customerWorkoutDay.UUID(), err)
		return ErrRepositoryFailure
	}

	group.UnregisterCustomer(args.CustomerUUID)
	err = a.repository.UpsertWorkoutGroup(ctx, group)
	if err != nil {
		const s = "upsert workout group UUID: %s workout day failed reason: %v"
		logger.Errorf(s, group.UUID(), err)
		return ErrRepositoryFailure
	}
	return nil
}

func NewUnassignCustomerHandler(c CustomerUnregister) *UnassignCustomerHandler {
	if c == nil {
		panic("nil repository")
	}
	return &UnassignCustomerHandler{
		repository: c,
	}
}
