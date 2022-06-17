package command

import (
	"context"
	"errors"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

type WorkoutScheduleDetails struct {
	CustomerUUID string
	GroupUUID    string
}

type ScheduleWorkoutHandlerRepository interface {
	QueryWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error)
	UpsertWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error
	UpsertCustomerWorkoutDay(ctx context.Context, workout customer.WorkoutDay) error
}

type ScheduleWorkoutHandler struct {
	repository ScheduleWorkoutHandlerRepository
}

func (s *ScheduleWorkoutHandler) Do(ctx context.Context, w WorkoutScheduleDetails) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "ScheduleWorkoutHandler"})
	group, err := s.repository.QueryWorkoutGroup(ctx, w.GroupUUID)
	if err != nil {
		logger.Errorf("query workout group UUID: %s failed, reason: %v", w.GroupUUID, err)
		return ErrRepositoryFailure
	}
	if group.UUID() != w.GroupUUID {
		logger.Errorf("group UUID: %s does not exist", w.GroupUUID)
		return ErrResourceNotFound
	}
	group.AssignCustomer(w.CustomerUUID)

	workoutDay, err := customer.NewWorkoutDay(w.CustomerUUID, group.UUID(), group.Date())
	if err != nil {
		const s = "creating workout day for customer UUID: %s, group UUID: %s, date: %s failed, reason: %v"
		logger.Errorf(s, w.CustomerUUID, w.GroupUUID, group.Date(), err)
		return err
	}

	err = s.repository.UpsertWorkoutGroup(ctx, group)
	if err != nil {
		logger.Errorf("upsert workout group UUID: %s failed, reason: %v", w.GroupUUID, err)
		return ErrRepositoryFailure
	}
	err = s.repository.UpsertCustomerWorkoutDay(ctx, *workoutDay)
	if err != nil {
		logger.Errorf("upsert customer workout day UUID: %s failed, reason: %v", workoutDay.UUID(), err)
		return ErrRepositoryFailure
	}
	return nil
}

func NewScheduleWorkoutHandler(c ScheduleWorkoutHandlerRepository) *ScheduleWorkoutHandler {
	if c == nil {
		panic("nil repository")
	}
	return &ScheduleWorkoutHandler{
		repository: c,
	}
}

var (
	ErrRepositoryFailure = errors.New("repository failure")
	ErrResourceNotFound  = errors.New("resource not found")
)
