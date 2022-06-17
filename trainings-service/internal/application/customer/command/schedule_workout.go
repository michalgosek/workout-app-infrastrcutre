package command

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/sirupsen/logrus"
)

type WorkoutScheduleDetails struct {
	CustomerUUID string
	GroupUUID    string
}

type ScheduleWorkoutHandler struct {
	trainerRepository  TrainerRepository
	customerRepository CustomerRepository
}

func (s *ScheduleWorkoutHandler) Do(ctx context.Context, w WorkoutScheduleDetails) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "ScheduleWorkoutHandler"})
	group, err := s.trainerRepository.QueryTrainerWorkoutGroup(ctx, w.GroupUUID)
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

	err = s.trainerRepository.UpsertTrainerWorkoutGroup(ctx, group)
	if err != nil {
		logger.Errorf("upsert workout group UUID: %s failed, reason: %v", w.GroupUUID, err)
		return ErrRepositoryFailure
	}
	err = s.customerRepository.UpsertCustomerWorkoutDay(ctx, *workoutDay)
	if err != nil {
		logger.Errorf("upsert customer workout day UUID: %s failed, reason: %v", workoutDay.UUID(), err)
		return ErrRepositoryFailure
	}
	return nil
}

func NewScheduleWorkoutHandler(c CustomerRepository, t TrainerRepository) *ScheduleWorkoutHandler {
	if c == nil {
		panic("nil customer repository")
	}
	if t == nil {
		panic("nil trainer repository")
	}
	return &ScheduleWorkoutHandler{
		customerRepository: c,
		trainerRepository:  t,
	}
}
