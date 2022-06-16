package command

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

type WorkoutRegistration struct {
	CustomerUUID string
	TrainerUUID  string
	GroupUUID    string
}

type CustomerAssigner interface {
	QueryWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error)
	UpsertWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error
	UpsertCustomerWorkoutDay(ctx context.Context, workout customer.WorkoutDay) error
}

type AssignCustomerHandler struct {
	repository CustomerAssigner
}

func (a *AssignCustomerHandler) Do(ctx context.Context, args WorkoutRegistration) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "AssignCustomerHandler"})

	group, err := a.repository.QueryWorkoutGroup(ctx, args.GroupUUID)
	if err != nil {
		return err
	}
	if group.TrainerUUID() != args.TrainerUUID {
		return ErrWorkoutGroupNotOwner
	}
	err = group.AssignCustomer(args.CustomerUUID)
	if err != nil {
		const s = "assign customer UUID: %s to the group UUID: %s failed, reason: %v"
		logger.Errorf(s, args.CustomerUUID, args.GroupUUID, err)
		return err
	}
	err = a.repository.UpsertWorkoutGroup(ctx, group)
	if err != nil {
		const s = "upsert customer UUID: %s to workout group UUID: %s failed, reason: %v"
		logger.Errorf(s, args.CustomerUUID, args.GroupUUID, err)
		return ErrRepositoryFailure
	}
	customerWorkoutDay, err := customer.NewWorkoutDay(args.CustomerUUID, args.GroupUUID, group.Date())
	if err != nil {
		const s = "creating new workout day for customer UUID: %s with workout group UUID: %s failed, reason: %v"
		logger.Errorf(s, args.CustomerUUID, args.GroupUUID, err)
		return err
	}
	err = a.repository.UpsertCustomerWorkoutDay(ctx, *customerWorkoutDay)
	if err != nil {
		const s = "upsert customer UUID: %s to workout day with group UUID: %s failed, reason: %v"
		logger.Errorf(s, args.CustomerUUID, args.GroupUUID, err)
		return ErrRepositoryFailure
	}
	return nil
}

func NewAssignCustomerHandler(c CustomerAssigner) *AssignCustomerHandler {
	if c == nil {
		panic("nil repository")
	}
	return &AssignCustomerHandler{
		repository: c,
	}
}
