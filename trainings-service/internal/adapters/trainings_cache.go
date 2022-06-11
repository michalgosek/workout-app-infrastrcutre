package adapters

import (
	"context"
	"errors"
	"fmt"
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


	Trainer:
	- Can schedule traninings (place limit, desc)
	- Can update training (place limit, desc, users)
	- Can cancel trainings
	- Can set trainigs limit (plcae limit) (Service layer -> UpdateWorkoutSessionLimit using under UpsertMethod if session exists)
	- Can remove user from the training (user)

	User (Member):
	- Can sing up to training -> Schedule
	- Can cancel his training
*/

type WorkoutsCacheRepoistory struct {
	workoutSessions     sync.Map
	userWorkoutSessions sync.Map
}

func NewWorkoutsCacheRepoistory() *WorkoutsCacheRepoistory {
	w := WorkoutsCacheRepoistory{}
	return &w
}

// TRAINER:
func (w *WorkoutsCacheRepoistory) UpsertTrainerWorkoutSessions(ctx context.Context, sessions ...domain.TrainerWorkoutSession) error {
	for _, s := range sessions {
		w.workoutSessions.Store(s.UUID, s)
	}
	return nil
}

func (w *WorkoutsCacheRepoistory) QueryTrainerWorkoutSessions(ctx context.Context, trainerUUID string) ([]domain.TrainerWorkoutSession, error) {
	var sessions []domain.TrainerWorkoutSession
	w.workoutSessions.Range(func(key, value interface{}) bool {
		session, ok := value.(domain.TrainerWorkoutSession)
		if !ok {
			return false
		}
		if session.TrainerUUID == trainerUUID {
			sessions = append(sessions, session)
		}
		return true
	})
	return sessions, nil
}

func (w *WorkoutsCacheRepoistory) QueryTrainerWorkoutSession(ctx context.Context, sessionUUID string) (domain.TrainerWorkoutSession, error) {
	session, ok := w.workoutSessions.Load(sessionUUID)
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
	session, ok := w.workoutSessions.Load(sessionUUID)
	if !ok {
		return domain.TrainerWorkoutSession{}, nil
	}
	v, ok := session.(domain.TrainerWorkoutSession)
	if !ok {
		return domain.TrainerWorkoutSession{}, fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, sessionUUID)
	}
	w.workoutSessions.Delete(sessionUUID)
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
func (w *WorkoutsCacheRepoistory) AddUserWorkoutSession(ctx context.Context, userUUID, sessionUUID string) error {
	return nil
}

func (w *WorkoutsCacheRepoistory) CancelUserWorkoutSession(ctx context.Context, userUUID, sessionUUID string) error {
	// reduce number of workouts in UserWorkout doc
	// reduce number of workouts in WorkoutGroup
	// increase number of user workotus limit
	return nil
}

func (w *WorkoutsCacheRepoistory) QueryUserWorkoutSessions(ctx context.Context, userUUID string) (domain.UserWorkoutSession, error) {
	return domain.UserWorkoutSession{}, nil
}

var ErrUnderlyingValueType = errors.New("invalid underlying value type")
