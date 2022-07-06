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
	logger := logrus.WithFields(logrus.Fields{"Component": "ScheduleWorkoutHandler", "Method": "checkCustomerWorkoutDuplicate"})
	workout, err := s.customerRepository.QueryCustomerWorkoutDay(ctx, w.CustomerUUID, w.GroupUUID)
	if err != nil {
		const s = "query customer UUID: %s  workout day with group UUID: %s failed, reason: %v"
		logger.Errorf(s, w.CustomerUUID, w.GroupUUID, err)
		return ErrRepositoryFailure
	}
	duplicate := workout.GroupUUID() == w.GroupUUID
	if duplicate {
		const s = "attempts to insert duplicate workout day with group UUID: %s for customer UUID: %s."
		logger.Warnf(s, w.GroupUUID, w.CustomerUUID)
		return ErrWorkoutGroupDuplicated
	}

	group, err := s.getTrainerWorkoutGroup(ctx, w)
	if err != nil {
		return err
	}
	err = s.scheduleWorkout(ctx, w, group)
	if err != nil {
		return err
	}
	return nil
}

func (s *ScheduleWorkoutHandler) getTrainerWorkoutGroup(ctx context.Context, w ScheduleWorkout) (trainer.WorkoutGroup, error) {
	logger := logrus.WithFields(logrus.Fields{"Component": "ScheduleWorkoutHandler", "Method": "getTrainerWorkoutGroup"})
	group, err := s.trainerRepository.QueryTrainerWorkoutGroup(ctx, w.TrainerUUID, w.GroupUUID)
	if err != nil {
		logger.Errorf("query workout group UUID: %s failed, reason: %v", w.GroupUUID, err)
		return trainer.WorkoutGroup{}, ErrRepositoryFailure
	}
	accessForbidden := group.UUID() != w.GroupUUID || group.TrainerUUID() != w.TrainerUUID
	if accessForbidden {
		logger.Warnf("group UUID: %s does not match to trainerUUID: %s.", w.GroupUUID, w.TrainerUUID)
		return trainer.WorkoutGroup{}, ErrResourceNotFound
	}
	return group, nil
}

func (s *ScheduleWorkoutHandler) scheduleWorkout(ctx context.Context, w ScheduleWorkout, trainerWorkout trainer.WorkoutGroup) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "ScheduleWorkoutHandler", "Method": "scheduleCustomerWorkout"})
	customerDetails, err := customer.NewCustomerDetails(w.CustomerUUID, w.CustomerName)
	if err != nil {
		logger.Errorf("creating customer details failed, reason: %v", err)
		return err
	}
	trainerWorkout.AssignCustomer(customerDetails)
	customerWorkout, err := customer.NewWorkoutDay(w.CustomerUUID, trainerWorkout.UUID(), trainerWorkout.Date())
	if err != nil {
		const s = "creating workout day for customer UUID: %s, group UUID: %s, date: %s failed, reason: %v"
		logger.Errorf(s, w.CustomerUUID, w.GroupUUID, trainerWorkout.Date(), err)
		return err
	}
	err = s.trainerRepository.UpsertTrainerWorkoutGroup(ctx, trainerWorkout)
	if err != nil {
		logger.Errorf("upsert workout group UUID: %s failed, reason: %v", w.GroupUUID, err)
		return ErrRepositoryFailure
	}
	err = s.customerRepository.UpsertCustomerWorkoutDay(ctx, customerWorkout)
	if err != nil {
		logger.Errorf("upsert customer workout day UUID: %s failed, reason: %v", customerWorkout.UUID(), err)
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
