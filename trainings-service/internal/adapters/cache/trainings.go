package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type TrainingSchedules struct {
	trainers  TrainerSchedules
	customers CustomerSchedules
}

func NewTrainingSchedules() *TrainingSchedules {
	w := TrainingSchedules{}
	return &w
}

func (w *TrainingSchedules) UpsertTrainerWorkoutSession(ctx context.Context, s trainer.TrainerSchedule) error {
	return w.trainers.UpsertSchedule(ctx, s)
}

func (w *TrainingSchedules) QueryTrainerWorkoutSessions(ctx context.Context, trainerUUID string) ([]trainer.TrainerSchedule, error) {
	return w.trainers.QuerySchedules(ctx, trainerUUID)
}

func (w *TrainingSchedules) QueryTrainerWorkoutSession(ctx context.Context, sessionUUID string) (trainer.TrainerSchedule, error) {
	return w.trainers.QuerySchedule(ctx, sessionUUID)
}

func (w *TrainingSchedules) CancelTrainerWorkoutSession(ctx context.Context, sessionUUID string) (trainer.TrainerSchedule, error) {
	return w.trainers.CancelSchedule(ctx, sessionUUID)
}

func (w *TrainingSchedules) CancelTrainerWorkoutSessions(ctx context.Context, sessionUUIDs ...string) ([]trainer.TrainerSchedule, error) {
	return w.trainers.CancelSchedules(ctx, sessionUUIDs...)
}

func (w *TrainingSchedules) UpsertCustomerWorkoutSession(ctx context.Context, session customer.CustomerSchedule) error {
	return w.customers.UpsertSchedule(ctx, session)
}

func (w *TrainingSchedules) UnregisterCustomerWorkoutSession(ctx context.Context, sessionUUID, customerUUID string) error {
	customerSession, err := w.customers.QuerySchedule(ctx, customerUUID)
	if err != nil {
		return fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, customerUUID)
	}

	trainerSession, err := w.trainers.QuerySchedule(ctx, sessionUUID)
	if err != nil {
		return fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, sessionUUID)
	}

	trainerSession.UnregisterCustomer(customerUUID)
	customerSession.UnregisterWorkout(sessionUUID)

	w.UpsertCustomerWorkoutSession(ctx, customerSession)
	w.UpsertTrainerWorkoutSession(ctx, trainerSession)
	return nil
}

func (w *TrainingSchedules) AssignCustomerToWorkoutSession(ctx context.Context, customerUUID, sessionUUID string) error {
	trainerWorkoutSession, err := w.QueryTrainerWorkoutSession(ctx, sessionUUID)
	if err != nil {
		return fmt.Errorf("query trainer workout session failed: %w", err)
	}
	customerWorkoutSession, err := w.QueryCustomerWorkoutSession(ctx, customerUUID)
	if err != nil {
		return fmt.Errorf("query customer workout session failed: %w", err)
	}
	if trainerWorkoutSession.UUID() == "" || customerWorkoutSession.UUID() == "" {
		return nil
	}

	err = trainerWorkoutSession.AssignCustomer(customerUUID)
	if err != nil {
		return fmt.Errorf("assign customer %s to workout session %s failed: %w", customerUUID, sessionUUID, err)
	}
	err = w.UpsertTrainerWorkoutSession(ctx, trainerWorkoutSession)
	if err != nil {
		return fmt.Errorf("upsert trainer workout session failed: %w", err)
	}

	err = customerWorkoutSession.AssignWorkout(sessionUUID)
	if err != nil {
		return fmt.Errorf("assign workout session %s to customer %s failed: %w", sessionUUID, customerUUID, err)
	}
	err = w.UpsertCustomerWorkoutSession(ctx, customerWorkoutSession)
	if err != nil {
		inner := w.UnregisterCustomerWorkoutSession(ctx, sessionUUID, customerUUID)
		if inner != nil {
			return fmt.Errorf("rollback failed: %w", err)
		}
		return fmt.Errorf("upsert trainer workout session failed: %w", err)
	}
	return nil
}

func (w *TrainingSchedules) QueryCustomerWorkoutSession(ctx context.Context, customerUUID string) (customer.CustomerSchedule, error) {
	return w.customers.QuerySchedule(ctx, customerUUID)
}

var ErrUnderlyingValueType = errors.New("invalid underlying value type")
