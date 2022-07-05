package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

type ScheduleWorkout struct {
	CustomerUUID string
	CustomerName string
	TrainerUUID  string
	GroupUUID    string
}

type ScheduleWorkoutHandler struct {
	trainerRepository  TrainerRepository
	customerRepository CustomerRepository
}

func (s *ScheduleWorkoutHandler) Do(ctx context.Context, w ScheduleWorkout) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "ScheduleWorkoutHandler"})

	duplicate, err := s.customerRepository.QueryCustomerWorkoutDay(ctx, w.CustomerUUID, w.GroupUUID)
	if err != nil {
		const s = "query customer UUID: %s  workout day with group UUID: %s failed, reason: %v"
		logger.Errorf(s, w.CustomerUUID, w.GroupUUID, err)
		return ErrRepositoryFailure
	}
	if duplicate.GroupUUID() == w.GroupUUID {
		const s = "attempts to insert duplicate workout day with group UUID: %s for customer UUID: %s failed, reason: %v"
		logger.Errorf(s, w.GroupUUID, w.CustomerUUID, err)
		return ErrWorkoutGroupDuplicated
	}

	group, err := s.trainerRepository.QueryTrainerWorkoutGroup(ctx, w.TrainerUUID, w.GroupUUID)
	if err != nil {
		logger.Errorf("query workout group UUID: %s failed, reason: %v", w.GroupUUID, err)
		return ErrRepositoryFailure
	}
	if group.UUID() != w.GroupUUID { //fixme also check match for trainer UUID!!
		logger.Errorf("group UUID: %s does not exist", w.GroupUUID)
		return ErrResourceNotFound
	}

	err = s.scheduleCustomerWorkout(ctx, w, group)
	if err != nil {
		return err
	}
	return nil
}

func (s *ScheduleWorkoutHandler) scheduleCustomerWorkout(ctx context.Context, w ScheduleWorkout, group trainer.WorkoutGroup) error {
	logger := logrus.WithFields(logrus.Fields{"Method": "scheduleCustomerWorkout"})
	customerDetails, err := customer.NewCustomerDetails(w.CustomerUUID, w.CustomerName)
	if err != nil {
		logger.Errorf("creating customer details failed, reason: %v", err)
		return err
	}
	group.AssignCustomer(customerDetails)
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
	err = s.customerRepository.UpsertCustomerWorkoutDay(ctx, workoutDay)
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
