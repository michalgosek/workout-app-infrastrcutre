package trainer

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Trainer:
// * Trainer cannot have more than 10 people during session and not less than 1
// * Training date must be not earlier than 3 hours from current date
// * Desc cannot be length than 100 chars
// * Name cannot be length than 15 chars

type TrainerSchedule struct {
	uuid          string
	trainerUUID   string
	limit         int
	customerUUIDs []string
	name          string
	desc          string
	date          time.Time
}

func (t *TrainerSchedule) Name() string {
	return t.name
}

func (t *TrainerSchedule) UUID() string {
	return t.uuid
}

func (t *TrainerSchedule) TrainerUUID() string {
	return t.trainerUUID
}

func (t *TrainerSchedule) Limit() int {
	return t.limit
}

func (t *TrainerSchedule) Date() time.Time {
	return t.date
}

func (t *TrainerSchedule) Desc() string {
	return t.desc
}

func (t *TrainerSchedule) CustomerUUIDs() []string {
	return t.customerUUIDs
}

func (t *TrainerSchedule) AssignedCustomers() int {
	return len(t.customerUUIDs)
}

func isProposedTimeNotExceeded(date time.Time) bool {
	threshold := time.Now().Add(3 * time.Hour)
	return date.Equal(threshold) || date.After(threshold)
}

func isProposedDescriptionNotExceeded(desc string) bool {
	return len(desc) > 100
}

func isProposedNameNotExceeded(name string) bool {
	return len(name) > 15
}

func (t *TrainerSchedule) UpdateDesc(s string) error {
	if isProposedDescriptionNotExceeded(s) {
		return ErrScheduleDescriptionExceeded
	}
	t.desc = s
	return nil
}

func (t *TrainerSchedule) UpdateName(s string) error {
	if isProposedNameNotExceeded(s) {
		return ErrScheduleNameExceeded
	}
	t.name = s
	return nil
}

func (t *TrainerSchedule) UpdateDate(d time.Time) error {
	if isProposedTimeNotExceeded(d) {
		return ErrScheduleDateViolation
	}
	t.date = d
	return nil
}

func (t *TrainerSchedule) UnregisterCustomer(UUID string) {
	var filtered []string
	for _, u := range t.customerUUIDs {
		if u == UUID {
			continue
		}
		filtered = append(filtered, u)
	}

	t.limit++
	t.customerUUIDs = filtered
}

func (t *TrainerSchedule) AssignCustomer(UUID string) error {
	if UUID == "" {
		return fmt.Errorf("%w: into customer workout session", ErrEmptyCustomerUUID)
	}
	if t.limit == 0 {
		return ErrCustomersScheduleLimitExceeded
	}
	if len(t.customerUUIDs) == 0 {
		t.customerUUIDs = append(t.customerUUIDs, UUID)
		t.limit--
		return nil
	}
	for _, u := range t.customerUUIDs {
		if u == UUID {
			return ErrDuplicateCustomerUUID
		}
	}
	t.customerUUIDs = append(t.customerUUIDs, UUID)
	t.limit--
	return nil
}

func NewSchedule(trainerUUID, name, desc string, date time.Time) (*TrainerSchedule, error) {
	ok := date.IsZero()
	if ok {
		return nil, ErrScheduleDateViolation
	}
	ok = isProposedTimeNotExceeded(date)
	if !ok {
		return nil, ErrScheduleDateViolation
	}
	ok = isProposedDescriptionNotExceeded(desc)
	if ok {
		return nil, ErrScheduleDescriptionExceeded
	}
	ok = isProposedNameNotExceeded(name)
	if ok {
		return nil, ErrScheduleNameExceeded
	}
	w := TrainerSchedule{
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

func UnmarshalFromDatabase(UUID, trainerUUID, name, desc string, customerUUIDs []string, date time.Time, limit int) TrainerSchedule {
	return TrainerSchedule{
		uuid:          UUID,
		trainerUUID:   trainerUUID,
		limit:         limit,
		customerUUIDs: customerUUIDs,
		name:          name,
		desc:          desc,
		date:          date,
	}
}

var (
	ErrCustomersScheduleLimitExceeded = errors.New("customers schedule limit exceeded")
	ErrScheduleNameExceeded           = errors.New("schedule name limit exceeded")
	ErrEmptyCustomerUUID              = errors.New("empty customer UUID ")
	ErrDuplicateCustomerUUID          = errors.New("customer UUID exists")
	ErrScheduleDateViolation          = errors.New("schedule date violation")
	ErrScheduleDescriptionExceeded    = errors.New("schedule description limit exceeded")
)
