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
	const (
		customerUUID            = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
		customerName            = "John Doe"
		trainerWorkoutGroupUUID = "01503798-eccb-4e90-8b12-d635e7494698"
	)
	date := time.Now().Add(24 * time.Hour)

	// when:
	actualWorkoutDay, err := customer.NewWorkoutDay(customerUUID, customerName, trainerWorkoutGroupUUID, date)

	// then:
	assertions.Nil(err)
	assertions.Equal(customerUUID, actualWorkoutDay.CustomerUUID())
	assertions.Equal(date, actualWorkoutDay.Date())
	assertions.Equal(trainerWorkoutGroupUUID, actualWorkoutDay.GroupUUID())
}

func TestCreateCustomerWorkoutDayShouldReturnErrorWhenSpecifiedEmptyCustomerName_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID            = "7e4946a1-8d33-4eec-b639-93fdb1c6fe30"
		customerName            = ""
		trainerWorkoutGroupUUID = "01503798-eccb-4e90-8b12-d635e7494698"
	)
	date := time.Now().Add(24 * time.Hour)

	// when:
	actualWorkoutDay, err := customer.NewWorkoutDay(customerUUID, customerName, trainerWorkoutGroupUUID, date)

	// then:
	assertions.Equal(err, customer.ErrEmptyCustomerName)
	assertions.Empty(actualWorkoutDay)
}

func TestCreateCustomerWorkoutDayShouldReturnErrorWhenSpecifiedEmptyCustomerUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	date := time.Now().Add(24 * time.Hour)
	const (
		trainerWorkoutGroupUUID = "01503798-eccb-4e90-8b12-d635e7494698"
		customerName            = "John Doe"
	)

	// when:
	actualWorkoutDay, err := customer.NewWorkoutDay("", customerName, trainerWorkoutGroupUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyCustomerUUID, err)
	assertions.Empty(actualWorkoutDay)
}

func TestCreateCustomerWorkoutDayShouldReturnErrorWhenSpecifiedEmptyDate_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID            = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
		customerName            = "John Doe"
		trainerWorkoutGroupUUID = "01503798-eccb-4e90-8b12-d635e7494698"
	)
	date := time.Time{}

	// when:
	actualWorkoutDay, err := customer.NewWorkoutDay(customerUUID, customerName, trainerWorkoutGroupUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyGroupDate, err)
	assertions.Empty(actualWorkoutDay)
}

func TestCreateCustomerWorkoutDayShouldReturnErrorWhenSpecifiedEmptyWorkoutUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID            = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
		customerName            = "John Doe"
		trainerWorkoutGroupUUID = ""
	)
	date := time.Now()

	// when:
	actualWorkoutDay, err := customer.NewWorkoutDay(customerUUID, customerName, trainerWorkoutGroupUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyGroupUUID, err)
	assertions.Empty(actualWorkoutDay)
}

func TestUnmarshalFromDatabaseShouldParseDataWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerWorkoutDayUUID  = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
		trainerWorkoutGroupUUID = "aa046e54-1c50-4b85-8a52-764d34c766ef"
		customerUUID            = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
		customerName            = "John Doe"
	)
	date := time.Now()

	// when:
	actualWorkoutDay, err := customer.UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID, customerName, date)

	// then:
	assertions.Nil(err)
	assertions.Equal(trainerWorkoutGroupUUID, actualWorkoutDay.GroupUUID())
	assertions.Equal(customerUUID, actualWorkoutDay.CustomerUUID())
	assertions.Equal(customerWorkoutDayUUID, actualWorkoutDay.UUID())
	assertions.Equal(customerName, actualWorkoutDay.CustomerName())
	assertions.Equal(date, actualWorkoutDay.Date())
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyCustomerName_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerWorkoutDayUUID  = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
		trainerWorkoutGroupUUID = "aa046e54-1c50-4b85-8a52-764d34c766ef"
		customerUUID            = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
		customerName            = ""
	)
	date := time.Now()

	// when:
	actualWorkoutDay, err := customer.UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID, customerName, date)

	// then:
	assertions.Equal(err, customer.ErrEmptyCustomerName)
	assertions.Empty(actualWorkoutDay)
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyCustomerWorkoutDayUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerWorkoutDayUUID  = ""
		trainerWorkoutGroupUUID = "aa046e54-1c50-4b85-8a52-764d34c766ef"
		customerUUID            = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
		customerName            = "John Doe"
	)
	date := time.Now()

	// when:
	SUT, err := customer.UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID, customerName, date)

	// then:
	assertions.Equal(customer.ErrEmptyWorkoutDayUUID, err)
	assertions.Empty(SUT)
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyCustomerUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerWorkoutDayUUID  = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
		customerUUID            = ""
		customerName            = "John Doe"
		trainerWorkoutGroupUUID = "aa046e54-1c50-4b85-8a52-764d34c766ef"
	)

	date := time.Now()

	// when:
	workout, err := customer.UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID, customerName, date)

	// then:
	assertions.Equal(customer.ErrEmptyCustomerUUID, err)
	assertions.Empty(workout)
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyGroupUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerWorkoutDayUUID  = "aa046e54-1c50-4b85-8a52-764d34c766ef"
		trainerWorkoutGroupUUID = ""
		customerName            = "John Doe"
		customerUUID            = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
	)
	date := time.Now()

	// when:
	SUT, err := customer.UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID, customerName, date)

	// then:
	assertions.Equal(customer.ErrEmptyGroupUUID, err)
	assertions.Empty(SUT)
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyGroupDate_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerWorkoutDayUUID  = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
		trainerWorkoutGroupUUID = "aa046e54-1c50-4b85-8a52-764d34c766ef"
		customerUUID            = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
		customerName            = "John Doe"
	)
	date := time.Time{}

	// when:
	SUT, err := customer.UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID, customerName, date)

	// then:
	assertions.Equal(customer.ErrEmptyGroupDate, err)
	assertions.Empty(SUT)
}
