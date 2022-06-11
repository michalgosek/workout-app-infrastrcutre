package domain

import (
	"time"

	"github.com/google/uuid"
)

// Optional
// Training can be canceled when users have been notified...

// Trainer:
// Trainer cannot have more than 10 people during session and not less than 1 (same applies to places)
// Training date must be not earlier than 3 hours from current date
// Desc cannot be length than 100 chars
// Name cannot be length than 15 chars
// places cannot be less than 0 or gerater than userUUIDS
// Trainer can create max 3 sessions

type TrainerWorkoutSession struct {
	uuid          string
	trainerUUID   string
	limit         int
	customerUUIDs []string
	name          string
	desc          string
	date          time.Time
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

func (t *TrainerWorkoutSession) AssignCustomer(UUIDs ...string) error {
	t.limit--
	t.customerUUIDs = append(t.customerUUIDs, UUIDs...)
	return nil
}

func NewTrainerWorkoutSession(trainerUUID, name, desc string, date time.Time) (*TrainerWorkoutSession, error) {
	// verify input
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
