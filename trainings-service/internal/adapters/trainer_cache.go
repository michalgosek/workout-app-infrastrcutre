package adapters

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain"
)

type TrainerSchedulesCache struct {
	lookup sync.Map
}

func (w *TrainerSchedulesCache) UpsertSchedule(ctx context.Context, s domain.TrainerWorkoutSession) error {
	w.lookup.Store(s.UUID(), s)
	return nil
}

func (t *TrainerSchedulesCache) QuerySchedules(_ context.Context, trainerUUID string) ([]domain.TrainerWorkoutSession, error) {
	var (
		sessions []domain.TrainerWorkoutSession
		keys     []string
	)
	t.lookup.Range(func(key, value interface{}) bool {
		session, ok := value.(domain.TrainerWorkoutSession)
		if !ok {
			return false
		}
		if session.TrainerUUID() == trainerUUID {
			keys = append(keys, session.UUID())
		}
		return true
	})

	sessions, err := t.sortSchedulesByUUID(keys, sessions)
	if err != nil {
		return nil, fmt.Errorf("sorting session by UUID failed: %w", err)
	}
	return sessions, nil
}

func (t *TrainerSchedulesCache) sortSchedulesByUUID(keys []string, sessions []domain.TrainerWorkoutSession) ([]domain.TrainerWorkoutSession, error) {
	sort.SliceStable(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, k := range keys {
		trainerSessionMapVal, ok := t.lookup.Load(k)
		if !ok {
			continue
		}
		session, ok := trainerSessionMapVal.(domain.TrainerWorkoutSession)
		if !ok {
			return nil, fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, k)
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (t *TrainerSchedulesCache) QuerySchedule(ctx context.Context, sessionUUID string) (domain.TrainerWorkoutSession, error) {
	trainerSessionMapVal, ok := t.lookup.Load(sessionUUID)
	if !ok {
		return domain.TrainerWorkoutSession{}, nil
	}
	session, ok := trainerSessionMapVal.(domain.TrainerWorkoutSession)
	if !ok {
		return domain.TrainerWorkoutSession{}, fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, sessionUUID)
	}
	return session, nil
}

func (w *TrainerSchedulesCache) CancelSchedule(ctx context.Context, sessionUUID string) (domain.TrainerWorkoutSession, error) {
	trainerSessionMapVal, ok := w.lookup.Load(sessionUUID)
	if !ok {
		return domain.TrainerWorkoutSession{}, nil
	}
	session, ok := trainerSessionMapVal.(domain.TrainerWorkoutSession)
	if !ok {
		return domain.TrainerWorkoutSession{}, fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, sessionUUID)
	}
	w.lookup.Delete(sessionUUID)
	return session, nil
}

func (t *TrainerSchedulesCache) CancelSchedules(ctx context.Context, sessionUUIDs ...string) ([]domain.TrainerWorkoutSession, error) {
	var sessions []domain.TrainerWorkoutSession
	for _, s := range sessionUUIDs {
		session, err := t.CancelSchedule(ctx, s)
		if err != nil {
			return nil, fmt.Errorf("delete trainer workout session failed: %w", err)
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func NewTrainerSchedulesCache() *TrainerSchedulesCache {
	return &TrainerSchedulesCache{}
}
