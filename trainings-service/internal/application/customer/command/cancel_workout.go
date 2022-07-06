package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

type CancelWorkout struct {
	CustomerUUID string
	TrainerUUID  string
	GroupUUID    string
}

type CancelWorkoutHandler struct {
	trainerRepository  TrainerRepository
	customerRepository CustomerRepository
}

func (c *CancelWorkoutHandler) Do(ctx context.Context, w CancelWorkout) error {
	group, err := c.getTrainerWorkoutGroup(ctx, w)
	if err != nil {
		return err
	}
	err = c.cancelWorkout(ctx, w, group)
	if err != nil {
		return err
	}
	return nil
}

func (c *CancelWorkoutHandler) cancelWorkout(ctx context.Context, w CancelWorkout, group trainer.WorkoutGroup) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "CancelWorkoutHandler", "Method": "cancelWorkout"})
	err := c.customerRepository.DeleteCustomerWorkoutDay(ctx, w.CustomerUUID, w.GroupUUID)
	if err != nil {
		logger.Errorf("delete customer UUID: %s workout day UUID: %s failed, reason: %v", w.CustomerUUID, w.GroupUUID, err)
		return ErrRepositoryFailure
	}
	group.UnregisterCustomer(w.CustomerUUID)
	err = c.trainerRepository.UpsertTrainerWorkoutGroup(ctx, group)
	if err != nil {
		logger.Errorf("upsert workout group UUID: %s failed, reason: %v", w.GroupUUID, err)
		return ErrRepositoryFailure
	}
	return nil
}

func (c *CancelWorkoutHandler) getTrainerWorkoutGroup(ctx context.Context, w CancelWorkout) (trainer.WorkoutGroup, error) {
	logger := logrus.WithFields(logrus.Fields{"Component": "CancelWorkoutHandler", "Method": "getTrainerWorkoutGroup"})
	workout, err := c.trainerRepository.QueryTrainerWorkoutGroup(ctx, w.TrainerUUID, w.GroupUUID)
	if err != nil {
		logger.Errorf("query workout group UUID: %s failed, reason: %v", w.GroupUUID, err)
		return trainer.WorkoutGroup{}, ErrRepositoryFailure
	}
	accessForbidden := workout.UUID() != w.GroupUUID || workout.TrainerUUID() != w.TrainerUUID
	if accessForbidden {
		logger.Warnf("group UUID: %s does not match to trainerUUID: %s.", w.GroupUUID, w.TrainerUUID)
		return trainer.WorkoutGroup{}, ErrResourceNotFound
	}
	return workout, nil
}

func NewCancelWorkoutHandler(c CustomerRepository, t TrainerRepository) *CancelWorkoutHandler {
	if c == nil {
		panic("nil customer repository")
	}
	if t == nil {
		panic("nil trainer repository")
	}
	return &CancelWorkoutHandler{trainerRepository: t, customerRepository: c}
}
