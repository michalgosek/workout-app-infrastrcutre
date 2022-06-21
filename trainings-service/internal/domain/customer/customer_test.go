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
	SUT, err := customer.NewWorkoutDay(customerUUID, trainerWorkoutGroupUUID, date)

	// then:
	assertions.Nil(err)
	assertions.Equal(customerUUID, SUT.CustomerUUID())
	assertions.Equal(date, SUT.Date())
	assertions.Equal(trainerWorkoutGroupUUID, SUT.GroupUUID())
}

func TestCreateCustomerWorkoutDayShouldReturnErrorWhenSpecifiedEmptyCustomerUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	date := time.Now().Add(24 * time.Hour)
	const trainerWorkoutGroupUUID = "01503798-eccb-4e90-8b12-d635e7494698"

	// when:
	SUT, err := customer.NewWorkoutDay("", trainerWorkoutGroupUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyCustomerUUID, err)
	assertions.Nil(SUT)
}

func TestCreateCustomerWorkoutDayShouldReturnErrorWhenSpecifiedEmptyDate_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const trainerWorkoutGroupUUID = "01503798-eccb-4e90-8b12-d635e7494698"
	date := time.Time{}

	// when:
	SUT, err := customer.NewWorkoutDay(customerUUID, trainerWorkoutGroupUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyGroupDate, err)
	assertions.Nil(SUT)
}

func TestCreateCustomerWorkoutDayShouldReturnErrorWhenSpecifiedEmptyWorkoutUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const trainerWorkoutGroupUUID = ""
	date := time.Now()

	// when:
	SUT, err := customer.NewWorkoutDay(customerUUID, trainerWorkoutGroupUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyGroupUUID, err)
	assertions.Nil(SUT)
}

func TestUnmarshalFromDatabaseShouldParseDataWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerWorkoutDayUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const trainerWorkoutGroupUUID = "aa046e54-1c50-4b85-8a52-764d34c766ef"
	const customerUUID = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
	date := time.Now()

	// when:
	SUT, err := customer.UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID, date)

	// then:
	assertions.Nil(err)
	assertions.Equal(trainerWorkoutGroupUUID, SUT.GroupUUID())
	assertions.Equal(customerUUID, SUT.CustomerUUID())
	assertions.Equal(customerWorkoutDayUUID, SUT.UUID())
	assertions.Equal(date, SUT.Date())
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyCustomerWorkoutDayUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerWorkoutDayUUID = ""
	const trainerWorkoutGroupUUID = "aa046e54-1c50-4b85-8a52-764d34c766ef"
	const customerUUID = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
	date := time.Now()

	// when:
	SUT, err := customer.UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyWorkoutDayUUID, err)
	assertions.Empty(SUT)
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
	SUT, err := customer.UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyGroupUUID, err)
	assertions.Empty(SUT)
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyGroupDate_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerWorkoutDayUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const trainerWorkoutGroupUUID = "aa046e54-1c50-4b85-8a52-764d34c766ef"
	const customerUUID = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
	date := time.Time{}

	// when:
	SUT, err := customer.UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyGroupDate, err)
	assertions.Empty(SUT)
}
