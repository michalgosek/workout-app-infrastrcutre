package trainer

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/verifiers"
)

/*
type UserSchedule struct {
	UserUUID        string
	Limit           int
	Schedules []string
}

type Schedule struct {
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
-  Can set trainigs limit (plcae limit) (Service layer -> UpdateScheduleLimit using under UpsertMethod if session exists)


	Trainer:
	- Can schedule traninings (place limit, desc)
	- Can update training (place limit, desc, users)
	- Can cancel trainings
	- Can remove user from the training (user)

*/

type TrainerSchedule struct {
	uuid          string
	trainerUUID   string
	limit         int
	customerUUIDs []string
	name          string
	desc          string
	date          time.Time
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

func (t *TrainerSchedule) Customers() int {
	return len(t.customerUUIDs)
}

func (t *TrainerSchedule) UpdateDesc(s string) error {
	t.desc = s
	return nil
}

func (t *TrainerSchedule) UpdateName(s string) error {
	t.name = s
	return nil
}

func (t *TrainerSchedule) SetDate(d time.Time) (time.Time, error) {
	t.date = d
	return t.date, nil
}

func (t *TrainerSchedule) UnregisterCustomer(UUID string) {
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

func (t *TrainerSchedule) AssignCustomer(UUID string) error {
	if UUID == "" {
		return fmt.Errorf("%w: into customer workout session", ErrEmptyCustomerUUID)
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
		return ErrCustomersScheduleLimitExceeded
	}

	t.customerUUIDs = append(t.customerUUIDs, UUID)
	t.limit--
	return nil
}

func NewSchedule(trainerUUID, name, desc string, date time.Time) (*TrainerSchedule, error) {
	dataVerifier := verifiers.NewWorkoutDate(3)
	err := dataVerifier.Check(date)
	if err != nil {
		return nil, err
	}
	descVerifier := verifiers.NewWorkoutDescription(100)
	err = descVerifier.Check(desc)
	if err != nil {
		return nil, err
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

var (
	ErrCustomersScheduleLimitExceeded = errors.New("customers schedule limit exceeded")
	ErrEmptyCustomerUUID              = errors.New("empty customer UUID ")
)
