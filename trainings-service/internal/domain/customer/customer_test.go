package customer_test

import (
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/stretchr/testify/assert"
)

func TestShouldCreateCustomerWorkoutDayWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const trainerWorkoutGroupUUID = "01503798-eccb-4e90-8b12-d635e7494698"
	date := time.Now().Add(24 * time.Hour)

	// when:
	workout, err := customer.NewWorkoutDay(customerUUID, trainerWorkoutGroupUUID, date)

	// then:
	assertions.Nil(err)
	assertions.Equal(customerUUID, workout.CustomerUUID())
	assertions.Equal(date, workout.Date())
	assertions.Equal(trainerWorkoutGroupUUID, workout.GroupUUID())
}

func TestCreateCustomerWorkoutDayShouldReturnErrorWhenSpecifiedEmptyCustomerUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	date := time.Now().Add(24 * time.Hour)
	const trainerWorkoutGroupUUID = "01503798-eccb-4e90-8b12-d635e7494698"

	// when:
	workout, err := customer.NewWorkoutDay("", trainerWorkoutGroupUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyCustomerUUID, err)
	assertions.Nil(workout)
}

func TestCreateCustomerWorkoutDayShouldReturnErrorWhenSpecifiedEmptyDate_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const trainerWorkoutGroupUUID = "01503798-eccb-4e90-8b12-d635e7494698"
	date := time.Time{}

	// when:
	workout, err := customer.NewWorkoutDay(customerUUID, trainerWorkoutGroupUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyGroupDate, err)
	assertions.Nil(workout)
}

func TestCreateCustomerWorkoutDayShouldReturnErrorWhenSpecifiedEmptyWorkoutUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const trainerWorkoutGroupUUID = ""
	date := time.Now()

	// when:
	workout, err := customer.NewWorkoutDay(customerUUID, trainerWorkoutGroupUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyGroupUUID, err)
	assertions.Nil(workout)
}

func TestUnmarshalFromDatabaseShouldParseDataWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerWorkoutDayUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const trainerWorkoutGroupUUID = "aa046e54-1c50-4b85-8a52-764d34c766ef"
	const customerUUID = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
	date := time.Now()

	// when:
	workout, err := customer.UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID, date)

	// then:
	assertions.Nil(err)
	assertions.Equal(trainerWorkoutGroupUUID, workout.GroupUUID())
	assertions.Equal(customerUUID, workout.CustomerUUID())
	assertions.Equal(customerWorkoutDayUUID, workout.UUID())
	assertions.Equal(date, workout.Date())
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyCustomerWorkoutDayUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerWorkoutDayUUID = ""
	const trainerWorkoutGroupUUID = "aa046e54-1c50-4b85-8a52-764d34c766ef"
	const customerUUID = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
	date := time.Now()

	// when:
	workout, err := customer.UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyWorkoutDayUUID, err)
	assertions.Empty(workout)
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyCustomerUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerWorkoutDayUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const trainerWorkoutGroupUUID = "aa046e54-1c50-4b85-8a52-764d34c766ef"
	const customerUUID = ""
	date := time.Now()

	// when:
	workout, err := customer.UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyCustomerUUID, err)
	assertions.Empty(workout)
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyGroupUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerWorkoutDayUUID = "aa046e54-1c50-4b85-8a52-764d34c766ef"
	const trainerWorkoutGroupUUID = ""
	const customerUUID = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
	date := time.Now()

	// when:
	workout, err := customer.UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyGroupUUID, err)
	assertions.Empty(workout)
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyGroupDate_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerWorkoutDayUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const trainerWorkoutGroupUUID = "aa046e54-1c50-4b85-8a52-764d34c766ef"
	const customerUUID = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
	date := time.Time{}

	// when:
	workout, err := customer.UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyGroupDate, err)
	assertions.Empty(workout)
}
