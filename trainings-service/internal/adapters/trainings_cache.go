package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain"
)

type WorkoutsCacheRepoistory struct {
	trainers  TrainerSchedulesCache
	customers CustomerSchedulesCache
}

func NewWorkoutsCacheRepoistory() *WorkoutsCacheRepoistory {
	w := WorkoutsCacheRepoistory{}
	return &w
}

func (w *WorkoutsCacheRepoistory) UpsertTrainerWorkoutSession(ctx context.Context, s domain.TrainerWorkoutSession) error {
	return w.trainers.UpsertSchedule(ctx, s)
}

func (w *WorkoutsCacheRepoistory) QueryTrainerWorkoutSessions(ctx context.Context, trainerUUID string) ([]domain.TrainerWorkoutSession, error) {
	return w.trainers.QuerySchedules(ctx, trainerUUID)
}

func (w *WorkoutsCacheRepoistory) QueryTrainerWorkoutSession(ctx context.Context, sessionUUID string) (domain.TrainerWorkoutSession, error) {
	return w.trainers.QuerySchedule(ctx, sessionUUID)
}

func (w *WorkoutsCacheRepoistory) CancelTrainerWorkoutSession(ctx context.Context, sessionUUID string) (domain.TrainerWorkoutSession, error) {
	return w.trainers.CancelSchedule(ctx, sessionUUID)
}

func (w *WorkoutsCacheRepoistory) CancelTrainerWorkoutSessions(ctx context.Context, sessionUUIDs ...string) ([]domain.TrainerWorkoutSession, error) {
	return w.trainers.CancelSchedules(ctx, sessionUUIDs...)
}

func (w *WorkoutsCacheRepoistory) UpsertCustomerWorkoutSession(ctx context.Context, session domain.CustomerWorkoutSession) error {
	return w.customers.UpsertSchedule(ctx, session)
}

func (w *WorkoutsCacheRepoistory) UnregisterCustomerWorkoutSession(ctx context.Context, sessionUUID, customerUUID string) error {
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

func (w *WorkoutsCacheRepoistory) AssignCustomerToWorkoutSession(ctx context.Context, customerUUID, sessionUUID string) error {
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

func (w *WorkoutsCacheRepoistory) QueryCustomerWorkoutSession(ctx context.Context, customerUUID string) (domain.CustomerWorkoutSession, error) {
	return w.customers.QuerySchedule(ctx, customerUUID)
}

var ErrUnderlyingValueType = errors.New("invalid underlying value type")
