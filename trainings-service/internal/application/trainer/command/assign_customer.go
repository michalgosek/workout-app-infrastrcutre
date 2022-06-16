package command

import (
	"context"
	"fmt"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type WorkoutRegistration struct {
	CustomerUUID string
	TrainerUUID  string
	GroupUUID    string
}

type CustomerAssigner interface {
	QueryWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error)
	UpsertWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error
}

type AssignCustomerHandler struct {
	repository CustomerAssigner
}

func (a *AssignCustomerHandler) Do(ctx context.Context, args WorkoutRegistration) error {
	group, err := a.repository.QueryWorkoutGroup(ctx, args.GroupUUID)
	if err != nil {
		return err
	}
	if group.TrainerUUID() != args.TrainerUUID {
		return ErrScheduleNotOwner
	}
	err = group.AssignCustomer(args.CustomerUUID)
	if err != nil {
		return fmt.Errorf("assign customer to the group failed: %w", err)
	}
	err = a.repository.UpsertWorkoutGroup(ctx, group)
	if err != nil {
		return fmt.Errorf("upsert group failed: %w", ErrRepositoryFailure)
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
