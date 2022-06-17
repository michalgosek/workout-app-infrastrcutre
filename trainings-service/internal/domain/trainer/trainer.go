package trainer

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Trainer:
// * Trainer cannot have more than 10 people during session and not less than 1
// * Training date must be not earlier than 3 hours from current date
// * Desc cannot be length than 100 chars
// * Name cannot be length than 15 chars

type WorkoutGroup struct {
	uuid          string
	trainerUUID   string
	trainerName   string
	limit         int
	customerUUIDs []string
	groupName     string
	groupDesc     string
	date          time.Time
}

func (t *WorkoutGroup) GroupName() string {
	return t.groupName
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

func (t *WorkoutGroup) GroupDescription() string {
	return t.groupDesc
}

func (t *WorkoutGroup) CustomerUUIDs() []string {
	return t.customerUUIDs
}

func (t *WorkoutGroup) AssignedCustomers() int {
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

func (t *WorkoutGroup) UpdateGroupDescription(s string) error {
	if isProposedDescriptionNotExceeded(s) {
		return ErrScheduleDescriptionExceeded
	}
	t.groupDesc = s
	return nil
}

func (t *WorkoutGroup) UpdateGroupName(s string) error {
	if isProposedNameNotExceeded(s) {
		return ErrScheduleNameExceeded
	}
	t.groupName = s
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

func (t *WorkoutGroup) AssignCustomer(UUID string) error {
	if UUID == "" {
		return ErrEmptyCustomerUUID
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
		uuid:          uuid.NewString(),
		trainerUUID:   trainerUUID,
		trainerName:   trainerName,
		groupName:     groupName,
		groupDesc:     groupDesc,
		limit:         10,
		date:          date,
		customerUUIDs: []string{},
	}
	return &w, nil
}

func UnmarshalFromDatabase(groupUUID, trainerUUID, trainerName, groupName, groupDesc string, customerUUIDs []string, date time.Time, limit int) (WorkoutGroup, error) {
	if groupUUID == "" {
		return WorkoutGroup{}, ErrEmptyWorkoutGroupUUID
	}
	if trainerUUID == "" {
		return WorkoutGroup{}, ErrEmptyWorkoutTrainerUUID
	}
	if trainerName == "" {
		return WorkoutGroup{}, ErrEmptyTrainerName
	}
	if groupName == "" {
		return WorkoutGroup{}, ErrEmptyWorkoutGroupName
	}
	if groupDesc == "" {
		return WorkoutGroup{}, ErrEmptyWorkoutGroupDesc
	}
	if date.IsZero() {
		return WorkoutGroup{}, ErrEmptyWorkoutGroupDate
	}
	w := WorkoutGroup{
		uuid:          groupUUID,
		trainerUUID:   trainerUUID,
		limit:         limit,
		customerUUIDs: customerUUIDs,
		groupName:     groupName,
		trainerName:   trainerName,
		groupDesc:     groupDesc,
		date:          date,
	}
	return w, nil
}

var (
	ErrEmptyWorkoutGroupUUID   = errors.New("empty workout group UUID")
	ErrEmptyWorkoutTrainerUUID = errors.New("empty trainer UUID")
	ErrEmptyWorkoutGroupName   = errors.New("empty workout name")
	ErrEmptyTrainerName        = errors.New("empty trainer name")
	ErrEmptyWorkoutGroupDesc   = errors.New("empty workout desc")
	ErrEmptyCustomerUUID       = errors.New("empty customer UUID")
	ErrEmptyWorkoutGroupDate   = errors.New("empty workout date")
)

var (
	ErrCustomersScheduleLimitExceeded = errors.New("customers schedule limit exceeded")
	ErrScheduleNameExceeded           = errors.New("schedule name limit exceeded")
	ErrDuplicateCustomerUUID          = errors.New("customer UUID exists")
	ErrScheduleDateViolation          = errors.New("schedule date violation")
	ErrScheduleDescriptionExceeded    = errors.New("schedule description limit exceeded")
)
