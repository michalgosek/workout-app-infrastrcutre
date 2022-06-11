package adapters

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain"
)

/*
type UserWorkoutSession struct {
	UserUUID        string
	Limit           int
	WorkoutSessions []string
}

type WorkoutSession struct {
	UUID        string
	TrainerUUID string
	Name        string
	Desc        string
	Places      int
	Canceled    bool
	Users       []string
	Date        time.Time
}

// Serwis:
// Uzytkownik nie powiniene sie zapisac na trening, gdy przekroczono limit (serwis layer)
-  Can set trainigs limit (plcae limit) (Service layer -> UpdateWorkoutSessionLimit using under UpsertMethod if session exists)


	Trainer:
	- Can schedule traninings (place limit, desc)
	- Can update training (place limit, desc, users)
	- Can cancel trainings
	- Can remove user from the training (user)

*/

type WorkoutsCacheRepoistory struct {
	trainerSessions  sync.Map
	customerSessions sync.Map
}

func NewWorkoutsCacheRepoistory() *WorkoutsCacheRepoistory {
	w := WorkoutsCacheRepoistory{}
	return &w
}

// TRAINER:
func (w *WorkoutsCacheRepoistory) UpsertTrainerWorkoutSession(ctx context.Context, s domain.TrainerWorkoutSession) error {
	w.trainerSessions.Store(s.UUID(), s)
	return nil
}

func (w *WorkoutsCacheRepoistory) RemoveCustomerFromTrainerWorkoutSession(ctx context.Context, sessionUUID, customerUUID string) error {
	v, ok := w.customerSessions.Load(customerUUID)
	if !ok {
		return nil
	}
	customerSession, ok := v.(domain.CustomerWorkoutSession)
	if !ok {
		return fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, customerUUID)
	}

	v, ok = w.trainerSessions.Load(sessionUUID)
	if !ok {
		return nil
	}
	trainerSession, ok := v.(domain.TrainerWorkoutSession)
	if !ok {
		return fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, sessionUUID)
	}

	trainerSession.UnregisterCustomer(customerUUID)
	customerSession.UnregisterWorkout(sessionUUID)

	w.UpsertCustomerWorkoutSession(ctx, customerSession)
	w.UpsertTrainerWorkoutSession(ctx, trainerSession)
	return nil
}

func (w *WorkoutsCacheRepoistory) QueryTrainerWorkoutSessions(ctx context.Context, trainerUUID string) ([]domain.TrainerWorkoutSession, error) {
	var (
		sessions []domain.TrainerWorkoutSession
		keys     []string
	)
	w.trainerSessions.Range(func(key, value interface{}) bool {
		session, ok := value.(domain.TrainerWorkoutSession)
		if !ok {
			return false
		}
		if session.TrainerUUID() == trainerUUID {
			keys = append(keys, session.UUID())
		}
		return true
	})

	sessions, err := w.sortTrainerSessionsByUUID(keys, sessions)
	if err != nil {
		return nil, fmt.Errorf("sorting session by UUID failed: %w", err)
	}
	return sessions, nil
}

func (w *WorkoutsCacheRepoistory) sortTrainerSessionsByUUID(keys []string, sessions []domain.TrainerWorkoutSession) ([]domain.TrainerWorkoutSession, error) {
	sort.SliceStable(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, k := range keys {
		v, ok := w.trainerSessions.Load(k)
		if !ok {
			continue
		}
		session, ok := v.(domain.TrainerWorkoutSession)
		if !ok {
			return nil, fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, k)
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (w *WorkoutsCacheRepoistory) QueryTrainerWorkoutSession(ctx context.Context, sessionUUID string) (domain.TrainerWorkoutSession, error) {
	session, ok := w.trainerSessions.Load(sessionUUID)
	if !ok {
		return domain.TrainerWorkoutSession{}, nil
	}
	v, ok := session.(domain.TrainerWorkoutSession)
	if !ok {
		return domain.TrainerWorkoutSession{}, fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, sessionUUID)
	}
	return v, nil
}

func (w *WorkoutsCacheRepoistory) DeleteTrainerWorkoutSession(ctx context.Context, sessionUUID string) (domain.TrainerWorkoutSession, error) {
	session, ok := w.trainerSessions.Load(sessionUUID)
	if !ok {
		return domain.TrainerWorkoutSession{}, nil
	}
	v, ok := session.(domain.TrainerWorkoutSession)
	if !ok {
		return domain.TrainerWorkoutSession{}, fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, sessionUUID)
	}
	w.trainerSessions.Delete(sessionUUID)
	return v, nil
}

func (w *WorkoutsCacheRepoistory) DeleteTrainerWorkoutSessions(ctx context.Context, sessionUUIDs ...string) ([]domain.TrainerWorkoutSession, error) {
	var sessions []domain.TrainerWorkoutSession
	for _, s := range sessionUUIDs {
		v, err := w.DeleteTrainerWorkoutSession(ctx, s)
		if err != nil {
			return nil, fmt.Errorf("delete trainer workout session failed: %w", err)
		}
		sessions = append(sessions, v)
	}
	return sessions, nil
}

// User:

// User (Member):
// - Can sing up to training -> Schedule
// - Can cancel his training

func (w *WorkoutsCacheRepoistory) UpsertCustomerWorkoutSession(ctx context.Context, s domain.CustomerWorkoutSession) error {
	w.customerSessions.Store(s.UserUUID(), s)
	return nil
}

func (w *WorkoutsCacheRepoistory) AssignCustomerToWorkoutSession(ctx context.Context, customerUUID, sessionUUID string) error {
	v, err := w.QueryTrainerWorkoutSession(ctx, sessionUUID)
	if err != nil {
		return fmt.Errorf("query trainer workout session failed: %w", err)
	}
	c, err := w.QueryCustomerWorkoutSession(ctx, customerUUID)
	if err != nil {
		return fmt.Errorf("query customer workout session failed: %w", err)
	}

	err = v.AssignCustomer(customerUUID)
	if err != nil {
		return fmt.Errorf("assign customer %s to workout session %s failed: %w", customerUUID, sessionUUID, err)
	}

	err = w.UpsertTrainerWorkoutSession(ctx, v)
	if err != nil {
		return fmt.Errorf("upsert trainer workout session failed: %w", err)
	}

	err = c.AssignWorkout(sessionUUID)
	if err != nil {
		return fmt.Errorf("assign workout session %s to customer %s failed: %w", sessionUUID, customerUUID, err)
	}
	err = w.UpsertCustomerWorkoutSession(ctx, c)
	if err != nil {
		inner := w.RemoveCustomerFromTrainerWorkoutSession(ctx, sessionUUID, customerUUID)
		if inner != nil {
			return fmt.Errorf("rollback failed: %w", err)
		}
		return fmt.Errorf("upsert trainer workout session failed: %w", err)
	}
	return nil
}

func (w *WorkoutsCacheRepoistory) DeleteCustomerWorkoutSession(ctx context.Context, userUUID, sessionUUID string) error {
	// reduce number of workouts in UserWorkout doc
	// reduce number of workouts in WorkoutGroup
	// increase number of user workotus limit

	return nil
}

func (w *WorkoutsCacheRepoistory) QueryCustomerWorkoutSession(ctx context.Context, customerUUID string) (domain.CustomerWorkoutSession, error) {
	session, ok := w.customerSessions.Load(customerUUID)
	if !ok {
		return domain.CustomerWorkoutSession{}, nil
	}
	v, ok := session.(domain.CustomerWorkoutSession)
	if !ok {
		return domain.CustomerWorkoutSession{}, fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, customerUUID)
	}
	return v, nil
}

var ErrUnderlyingValueType = errors.New("invalid underlying value type")
