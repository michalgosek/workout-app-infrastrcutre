package cache

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type TrainerSchedules struct {
	lookup sync.Map
}

func (t *TrainerSchedules) UpsertSchedule(ctx context.Context, schedule trainer.TrainerSchedule) error {
	t.lookup.Store(schedule.UUID(), schedule)
	return nil
}

func (t *TrainerSchedules) QuerySchedules(_ context.Context, trainerUUID string) ([]trainer.TrainerSchedule, error) {
	var (
		schedules []trainer.TrainerSchedule
		keys      []string
	)
	t.lookup.Range(func(key, value interface{}) bool {
		schedule, ok := value.(trainer.TrainerSchedule)
		if !ok {
			return false
		}
		if schedule.TrainerUUID() == trainerUUID {
			keys = append(keys, schedule.UUID())
		}
		return true
	})
	schedules, err := t.sortSchedulesByUUID(keys)
	if err != nil {
		return nil, fmt.Errorf("sorting session by UUID failed: %w", err)
	}
	return schedules, nil
}

func (t *TrainerSchedules) sortSchedulesByUUID(keys []string) ([]trainer.TrainerSchedule, error) {
	var schedules []trainer.TrainerSchedule
	sort.SliceStable(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, k := range keys {
		trainerScheduleMapV, ok := t.lookup.Load(k)
		if !ok {
			continue
		}
		schedule, ok := trainerScheduleMapV.(trainer.TrainerSchedule)
		if !ok {
			return nil, fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, k)
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

func (t *TrainerSchedules) QuerySchedule(ctx context.Context, scheduleUUID string) (trainer.TrainerSchedule, error) {
	trainerScheduleMapV, ok := t.lookup.Load(scheduleUUID)
	if !ok {
		return trainer.TrainerSchedule{}, nil
	}
	schedule, ok := trainerScheduleMapV.(trainer.TrainerSchedule)
	if !ok {
		return trainer.TrainerSchedule{}, fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, scheduleUUID)
	}
	return schedule, nil
}

func (w *TrainerSchedules) CancelSchedule(ctx context.Context, sessionUUID string) (trainer.TrainerSchedule, error) {
	trainerScheduleMapV, ok := w.lookup.Load(sessionUUID)
	if !ok {
		return trainer.TrainerSchedule{}, nil
	}
	schedule, ok := trainerScheduleMapV.(trainer.TrainerSchedule)
	if !ok {
		return trainer.TrainerSchedule{}, fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, sessionUUID)
	}
	w.lookup.Delete(sessionUUID)
	return schedule, nil
}

func (t *TrainerSchedules) CancelSchedules(ctx context.Context, sessionUUIDs ...string) ([]trainer.TrainerSchedule, error) {
	var schedules []trainer.TrainerSchedule
	for _, s := range sessionUUIDs {
		schedule, err := t.CancelSchedule(ctx, s)
		if err != nil {
			return nil, fmt.Errorf("cancel trainer workout schedule failed: %w", err)
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

func newTrainerSchedules() *TrainerSchedules {
	return &TrainerSchedules{}
}
