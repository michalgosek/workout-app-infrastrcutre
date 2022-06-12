package domain

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/aggregates"
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

type TrainerWorkoutSession struct {
	uuid          string
	trainerUUID   string
	limit         int
	customerUUIDs []string
	name          string
	desc          string
	date          time.Time
}

func (t *TrainerWorkoutSession) AssignedCustomers() int {
	return len(t.customerUUIDs)
}

func (t *TrainerWorkoutSession) UUID() string {
	return t.uuid
}

func (t *TrainerWorkoutSession) TrainerUUID() string {
	return t.trainerUUID
}

func (t *TrainerWorkoutSession) Limit() int {
	return t.limit
}

func (t *TrainerWorkoutSession) Date() time.Time {
	return t.date
}

func (t *TrainerWorkoutSession) Customers() int {
	return len(t.customerUUIDs)
}

func (t *TrainerWorkoutSession) SetDesc(s string) error {
	t.desc = s
	return nil
}

func (t *TrainerWorkoutSession) SetName(s string) error {
	t.name = s
	return nil
}

func (t *TrainerWorkoutSession) SetDate(d time.Time) (time.Time, error) {
	t.date = d
	return t.date, nil
}

func (t *TrainerWorkoutSession) UnregisterCustomer(UUID string) {
	var filtered []string
	for _, u := range t.customerUUIDs {
		if u == UUID {
			continue
		}
		filtered = append(filtered, u)
	}

	t.limit = len(filtered)
	t.customerUUIDs = filtered
}

func (t *TrainerWorkoutSession) AssignCustomer(UUID string) error {
	if UUID == "" {
		return fmt.Errorf("%w: into customer workout session", ErrEmptyCustomerWorkoutSessionUUID)
	}
	if len(t.customerUUIDs) == 0 {
		t.customerUUIDs = append(t.customerUUIDs, UUID)
		t.limit--
		return nil
	}

	sort.SliceStable(t.customerUUIDs, func(i, j int) bool { return t.customerUUIDs[i] < t.customerUUIDs[j] })
	idx := sort.Search(len(t.customerUUIDs), func(i int) bool { return t.customerUUIDs[i] == UUID })
	if idx < len(t.customerUUIDs) && t.customerUUIDs[idx] == UUID {
		return nil
	}
	if t.limit == 0 {
		return ErrCustomerWorkouSessionLimitExceeded
	}

	t.customerUUIDs = append(t.customerUUIDs, UUID)
	t.limit--
	return nil
}

func NewTrainerWorkoutSession(trainerUUID, name, desc string, date time.Time) (*TrainerWorkoutSession, error) {
	dateAggregate := aggregates.NewWorkoutDate(3)
	err := dateAggregate.Check(date)
	if err != nil {
		return nil, err
	}
	descAggregate := aggregates.NewWorkoutDescription(100)
	err = descAggregate.Check(desc)
	if err != nil {
		return nil, err
	}
	w := TrainerWorkoutSession{
		uuid:          uuid.NewString(),
		trainerUUID:   trainerUUID,
		name:          name,
		desc:          desc,
		limit:         10,
		date:          date,
		customerUUIDs: []string{},
	}
	return &w, nil
}

var (
	ErrTrainerWorkouSessionLimitExceeded = errors.New("customer session workouts number exceeded")
	ErrEmptyCustomerWorkoutSessionUUID   = errors.New("empty customer workout session UUID")
)
