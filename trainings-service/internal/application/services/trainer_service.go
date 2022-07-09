package services

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"time"
)

//go:generate mockery --name=TrainerRepository --case underscore --with-expecter
type TrainerRepository interface {
	UpsertTrainerWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error
	QueryTrainerWorkoutGroup(ctx context.Context, trainerUUID, groupUUID string) (trainer.WorkoutGroup, error)
}

type AssignCustomerToWorkoutGroupArgs struct {
	CustomerUUID string
	CustomerName string
	GroupUUID    string
	TrainerUUID  string
}

type CancelCustomerWorkoutParticipationArgs struct {
	CustomerUUID string
	GroupUUID    string
	TrainerUUID  string
}

type TrainerService struct {
	repository TrainerRepository
}

type GetTrainerWorkoutGroupArgs struct {
	TrainerUUID string
	GroupUUID   string
}

type AssignedCustomerWorkoutGroupDetails struct {
	UUID        string
	TrainerUUID string
	Name        string
	Date        time.Time
}

func (t *TrainerService) AssignCustomerToWorkoutGroup(ctx context.Context, args AssignCustomerToWorkoutGroupArgs) (AssignedCustomerWorkoutGroupDetails, error) {
	group, err := t.repository.QueryTrainerWorkoutGroup(ctx, args.TrainerUUID, args.GroupUUID)
	if err != nil {
		return AssignedCustomerWorkoutGroupDetails{}, ErrQueryTrainerWorkoutGroup
	}
	if group.UUID() == "" {
		return AssignedCustomerWorkoutGroupDetails{}, ErrResourceNotFound
	}
	if group.UUID() == args.GroupUUID {
		return AssignedCustomerWorkoutGroupDetails{}, ErrResourceDuplicated
	}

	customerDetails, err := customer.NewCustomerDetails(args.CustomerUUID, args.CustomerName)
	if err != nil {
		return AssignedCustomerWorkoutGroupDetails{}, err
	}
	group.AssignCustomer(customerDetails)
	err = t.repository.UpsertTrainerWorkoutGroup(ctx, group)
	if err != nil {
		return AssignedCustomerWorkoutGroupDetails{}, ErrUpsertTrainerWorkoutGroup
	}
	workoutGroupDetails := AssignedCustomerWorkoutGroupDetails{
		Date:        group.Date(),
		UUID:        group.UUID(),
		TrainerUUID: group.TrainerUUID(),
		Name:        group.Name(),
	}
	return workoutGroupDetails, nil
}

func (t *TrainerService) CancelCustomerWorkoutParticipation(ctx context.Context, args CancelCustomerWorkoutParticipationArgs) error {
	group, err := t.repository.QueryTrainerWorkoutGroup(ctx, args.TrainerUUID, args.GroupUUID)
	if err != nil {
		return ErrQueryTrainerWorkoutGroup
	}
	// - group empty means there is no group for trainer

	group.UnregisterCustomer(args.CustomerUUID)
	err = t.repository.UpsertTrainerWorkoutGroup(ctx, group)
	if err != nil {
		return ErrUpsertTrainerWorkoutGroup
	}
	return nil
}

func NewTrainerService(t TrainerRepository) (*TrainerService, error) {
	if t == nil {
		return nil, ErrNilTrainerRepository
	}
	s := TrainerService{repository: t}
	return &s, nil
}

var (
	ErrNilTrainerRepository = errors.New("nil trainer repository dependency")
)

var (
	ErrUpsertTrainerWorkoutGroup = errors.New("cmd upsert trainer workout group failure")
	ErrQueryTrainerWorkoutGroup  = errors.New("query trainer group failure")
)
