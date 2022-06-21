package trainer

import (
	"errors"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"

	"github.com/google/uuid"
)

// Trainer:
// * Trainer cannot have more than 10 people during session and not less than 1
// * Training date must be not earlier than 3 hours from current date
// * Desc cannot be length than 100 chars
// * Name cannot be length than 15 chars

type WorkoutGroup struct {
	uuid            string
	trainerUUID     string
	trainerName     string
	limit           int
	customerDetails []customer.Details
	name            string
	description     string
	date            time.Time
}

func (t *WorkoutGroup) Name() string {
	return t.name
}

func (t *WorkoutGroup) TrainerName() string {
	return t.trainerName
}

func (t *WorkoutGroup) UUID() string {
	return t.uuid
}

func (t *WorkoutGroup) TrainerUUID() string {
	return t.trainerUUID
}

func (t *WorkoutGroup) Limit() int {
	return t.limit
}

func (t *WorkoutGroup) Date() time.Time {
	return t.date
}

func (t *WorkoutGroup) Description() string {
	return t.description
}

func (t *WorkoutGroup) CustomerDetails() []customer.Details {
	return t.customerDetails
}

func (t *WorkoutGroup) AssignedCustomers() int {
	return len(t.customerDetails)
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

func (t *WorkoutGroup) UpdateDescription(s string) error {
	if isProposedDescriptionNotExceeded(s) {
		return ErrScheduleDescriptionExceeded
	}
	t.description = s
	return nil
}

func (t *WorkoutGroup) UpdateName(s string) error {
	if isProposedNameNotExceeded(s) {
		return ErrScheduleNameExceeded
	}
	t.name = s
	return nil
}

func (t *WorkoutGroup) UpdateGroupDate(d time.Time) error {
	if isProposedTimeNotExceeded(d) {
		return ErrScheduleDateViolation
	}
	t.date = d
	return nil
}

func (t *WorkoutGroup) UnregisterCustomer(UUID string) {
	var filtered []customer.Details
	for _, c := range t.CustomerDetails() {
		if c.UUID() == UUID {
			continue
		}
		filtered = append(filtered, c)
	}
	t.limit++
	t.customerDetails = filtered
}

func (t *WorkoutGroup) AssignCustomer(c customer.Details) error {
	if c.UUID() == "" {
		return customer.ErrEmptyCustomerUUID
	}
	if t.limit == 0 {
		return ErrCustomersScheduleLimitExceeded
	}
	if len(t.customerDetails) == 0 {
		t.customerDetails = append(t.customerDetails, c)
		t.limit--
		return nil
	}
	for _, d := range t.customerDetails {
		if d.UUID() == c.UUID() {
			return ErrDuplicateCustomerUUID
		}
	}
	t.customerDetails = append(t.customerDetails, c)
	t.limit--
	return nil
}

func NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc string, date time.Time) (*WorkoutGroup, error) {
	ok := date.IsZero()
	if ok {
		return nil, ErrScheduleDateViolation
	}
	ok = isProposedTimeNotExceeded(date)
	if !ok {
		return nil, ErrScheduleDateViolation
	}
	ok = isProposedDescriptionNotExceeded(groupDesc)
	if ok {
		return nil, ErrScheduleDescriptionExceeded
	}
	ok = isProposedNameNotExceeded(groupName)
	if ok {
		return nil, ErrScheduleNameExceeded
	}
	if trainerName == "" {
		return nil, ErrEmptyTrainerName
	}
	w := WorkoutGroup{
		uuid:        uuid.NewString(),
		trainerUUID: trainerUUID,
		trainerName: trainerName,
		name:        groupName,
		description: groupDesc,
		limit:       10,
		date:        date,
	}
	return &w, nil
}

type WorkoutGroupDetails struct {
	UUID        string
	TrainerUUID string
	TrainerName string
	Name        string
	Description string
	Date        time.Time
	Limit       int
}

func UnmarshalWorkoutGroupFromDatabase(w WorkoutGroupDetails, customerDetails []customer.Details) (WorkoutGroup, error) {
	if w.UUID == "" {
		return WorkoutGroup{}, ErrEmptyWorkoutGroupUUID
	}
	if w.Name == "" {
		return WorkoutGroup{}, ErrEmptyWorkoutGroupName
	}
	if w.Description == "" {
		return WorkoutGroup{}, ErrEmptyWorkoutGroupDesc
	}
	if w.Date.IsZero() {
		return WorkoutGroup{}, ErrEmptyWorkoutGroupDate
	}
	if w.TrainerUUID == "" {
		return WorkoutGroup{}, ErrEmptyWorkoutTrainerUUID
	}
	if w.TrainerName == "" {
		return WorkoutGroup{}, ErrEmptyTrainerName
	}
	group := WorkoutGroup{
		uuid:            w.UUID,
		limit:           w.Limit,
		name:            w.Name,
		description:     w.Description,
		date:            w.Date,
		customerDetails: customerDetails,
		trainerUUID:     w.TrainerUUID,
		trainerName:     w.TrainerName,
	}
	return group, nil
}

var (
	ErrEmptyWorkoutGroupUUID   = errors.New("empty workout group UUID")
	ErrEmptyWorkoutTrainerUUID = errors.New("empty trainer UUID")
	ErrEmptyWorkoutGroupName   = errors.New("empty workout name")
	ErrEmptyTrainerName        = errors.New("empty trainer name")
	ErrEmptyWorkoutGroupDesc   = errors.New("empty workout desc")
	ErrEmptyWorkoutGroupDate   = errors.New("empty workout date")
)

var (
	ErrCustomersScheduleLimitExceeded = errors.New("customers schedule limit exceeded")
	ErrScheduleNameExceeded           = errors.New("schedule name limit exceeded")
	ErrDuplicateCustomerUUID          = errors.New("customer UUID exists")
	ErrScheduleDateViolation          = errors.New("schedule date violation")
	ErrScheduleDescriptionExceeded    = errors.New("schedule description limit exceeded")
)
