package trainer

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"time"
)

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

type CancelWorkoutGroupArgs struct {
	TrainerUUID string
	GroupUUID   string
}

type WorkoutGroupWithCustomerArgs struct {
	TrainerUUID  string
	GroupUUID    string
	CustomerUUID string
}

type CreateWorkoutGroupArgs struct {
	TrainerUUID string
	TrainerName string
	GroupName   string
	GroupDesc   string
	Date        time.Time
}

type Commands interface {
	UpsertTrainerWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error
	DeleteTrainerWorkoutGroup(ctx context.Context, trainerUUID, groupUUID string) error
	DeleteTrainerWorkoutGroups(ctx context.Context, trainerUUID string) error
}

type Queries interface {
	QueryTrainerWorkoutGroup(ctx context.Context, trainerUUID, groupUUID string) (trainer.WorkoutGroup, error)
	QueryTrainerWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error)
	QueryCustomerWorkoutGroup(ctx context.Context, trainerUUID, groupUUID, customerUUID string) (trainer.WorkoutGroup, error)
	QueryTrainerWorkoutGroupWithDate(ctx context.Context, trainerUUID string, date time.Time) (trainer.WorkoutGroup, error)
}

//go:generate mockery --name=Repository --case underscore --with-expecter
type Repository interface {
	Commands
	Queries
}

type Service struct {
	repository Repository
}

func (s *Service) GetTrainerWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error) {
	groups, err := s.repository.QueryTrainerWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		return nil, ErrErrQueryTrainerWorkoutGroups
	}
	return groups, nil
}

func (s *Service) CreateWorkoutGroup(ctx context.Context, args CreateWorkoutGroupArgs) error {
	duplicate, err := s.repository.QueryTrainerWorkoutGroupWithDate(ctx, args.TrainerUUID, args.Date)
	if err != nil {
		return ErrQueryTrainerWorkoutGroupWithDateWithDate
	}
	if duplicate.UUID() != "" {
		return ErrResourceDuplicated
	}
	group, err := trainer.NewWorkoutGroup(args.TrainerUUID, args.TrainerName, args.GroupName, args.GroupDesc, args.Date)
	if err != nil {
		return err
	}
	err = s.repository.UpsertTrainerWorkoutGroup(ctx, group)
	if err != nil {
		return ErrUpsertTrainerWorkoutGroup
	}
	return nil
}

func (s *Service) CancelWorkoutGroup(ctx context.Context, args CancelWorkoutGroupArgs) error {
	group, err := s.repository.QueryTrainerWorkoutGroup(ctx, args.TrainerUUID, args.GroupUUID)
	if err != nil {
		return ErrQueryTrainerWorkoutGroup
	}
	if group.UUID() != args.GroupUUID || group.TrainerUUID() != args.TrainerUUID {
		return ErrResourceNotFound
	}

	err = s.repository.DeleteTrainerWorkoutGroup(ctx, args.TrainerUUID, args.GroupUUID)
	if err != nil {
		return ErrDeleteTrainerWorkoutGroup
	}
	return nil
}

func (s *Service) CancelWorkoutGroups(ctx context.Context, trainerUUID string) error {
	err := s.repository.DeleteTrainerWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		return ErrDeleteTrainerWorkoutGroups
	}
	return nil
}

func (s *Service) AssignCustomerToWorkoutGroup(ctx context.Context, args AssignCustomerToWorkoutGroupArgs) (AssignedCustomerWorkoutGroupDetails, error) {
	group, err := s.repository.QueryTrainerWorkoutGroup(ctx, args.TrainerUUID, args.GroupUUID) // 	QueryTrainerWorkoutGroupWithCustomer()
	if err != nil {
		return AssignedCustomerWorkoutGroupDetails{}, ErrQueryTrainerWorkoutGroup
	}
	if group.UUID() == "" {
		return AssignedCustomerWorkoutGroupDetails{}, ErrResourceNotFound
	}

	customerDetails, err := customer.NewCustomerDetails(args.CustomerUUID, args.CustomerName)
	if err != nil {
		return AssignedCustomerWorkoutGroupDetails{}, err
	}
	group.AssignCustomer(customerDetails)
	err = s.repository.UpsertTrainerWorkoutGroup(ctx, group)
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

func (s *Service) CancelCustomerWorkoutParticipation(ctx context.Context, args CancelCustomerWorkoutParticipationArgs) error {
	group, err := s.repository.QueryTrainerWorkoutGroup(ctx, args.TrainerUUID, args.GroupUUID)
	if err != nil {
		return ErrQueryTrainerWorkoutGroup
	}
	// - group empty means there is no group for trainer

	group.UnregisterCustomer(args.CustomerUUID)
	err = s.repository.UpsertTrainerWorkoutGroup(ctx, group)
	if err != nil {
		return ErrUpsertTrainerWorkoutGroup
	}
	return nil
}

func (s *Service) GetWorkoutGroupWithCustomer(ctx context.Context, args WorkoutGroupWithCustomerArgs) (trainer.WorkoutGroup, error) {
	group, err := s.repository.QueryCustomerWorkoutGroup(ctx, args.TrainerUUID, args.GroupUUID, args.CustomerUUID)
	if err != nil {
		return trainer.WorkoutGroup{}, ErrQueryTrainerWorkoutGroupWithCustomer
	}
	return group, nil
}

func NewTrainerService(r Repository) (*Service, error) {
	if r == nil {
		return nil, ErrNilTrainerRepository
	}
	s := Service{repository: r}
	return &s, nil
}

var (
	ErrNilTrainerRepository = errors.New("nil trainer repository dependency")
)

var (
	ErrResourceNotFound   = errors.New("resource not found")
	ErrResourceDuplicated = errors.New("resource duplicated")
)

var (
	ErrDeleteTrainerWorkoutGroup                = errors.New("cmd delete trainer workout group failure")
	ErrDeleteTrainerWorkoutGroups               = errors.New("cmd delete trainer workout groups failure")
	ErrUpsertTrainerWorkoutGroup                = errors.New("cmd upsert trainer workout group failure")
	ErrQueryTrainerWorkoutGroup                 = errors.New("query trainer group failure")
	ErrErrQueryTrainerWorkoutGroups             = errors.New("query trainer workout groups failure")
	ErrQueryTrainerWorkoutGroupWithCustomer     = errors.New("query trainer workout group with customer failure")
	ErrQueryTrainerWorkoutGroupWithDateWithDate = errors.New("query trainer workout group with date failure")
)
