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
		customerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
		customerName = "John Doe"
		trainerUUID  = "d1f9cf76-0acf-445c-b9ee-d9515e20d22d"
		groupUUID    = "01503798-eccb-4e90-8b12-d635e7494698"
	)
	date := time.Now().Add(24 * time.Hour)

	// when:
	actualWorkoutDay, err := customer.NewWorkoutDay(customerUUID, customerName, groupUUID, trainerUUID, date)

	// then:
	assertions.Nil(err)
	assertions.Equal(customerUUID, actualWorkoutDay.CustomerUUID())
	assertions.Equal(date, actualWorkoutDay.Date())
	assertions.Equal(groupUUID, actualWorkoutDay.GroupUUID())
}

func TestCreateCustomerWorkoutDayShouldReturnErrorWhenSpecifiedEmptyCustomerName_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "7e4946a1-8d33-4eec-b639-93fdb1c6fe30"
		customerName = ""
		trainerUUID  = "d1f9cf76-0acf-445c-b9ee-d9515e20d22d"
		groupUUID    = "01503798-eccb-4e90-8b12-d635e7494698"
	)
	date := time.Now().Add(24 * time.Hour)

	// when:
	actualWorkoutDay, err := customer.NewWorkoutDay(customerUUID, customerName, trainerUUID, groupUUID, date)

	// then:
	assertions.Equal(err, customer.ErrEmptyCustomerName)
	assertions.Empty(actualWorkoutDay)
}

func TestCreateCustomerWorkoutDayShouldReturnErrorWhenSpecifiedEmptyCustomerUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	date := time.Now().Add(24 * time.Hour)
	const (
		groupUUID    = "01503798-eccb-4e90-8b12-d635e7494698"
		trainerUUID  = "d1f9cf76-0acf-445c-b9ee-d9515e20d22d"
		customerName = "John Doe"
	)

	// when:
	actualWorkoutDay, err := customer.NewWorkoutDay("", customerName, trainerUUID, groupUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyCustomerUUID, err)
	assertions.Empty(actualWorkoutDay)
}

func TestCreateCustomerWorkoutDayShouldReturnErrorWhenSpecifiedEmptyDate_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
		customerName = "John Doe"
		trainerUUID  = "d1f9cf76-0acf-445c-b9ee-d9515e20d22d"
		groupUUID    = "01503798-eccb-4e90-8b12-d635e7494698"
	)
	date := time.Time{}

	// when:
	actualWorkoutDay, err := customer.NewWorkoutDay(customerUUID, customerName, trainerUUID, groupUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyGroupDate, err)
	assertions.Empty(actualWorkoutDay)
}

func TestCreateCustomerWorkoutDayShouldReturnErrorWhenSpecifiedEmptyWorkoutUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
		customerName = "John Doe"
		trainerUUID  = "d1f9cf76-0acf-445c-b9ee-d9515e20d22d"
		groupUUID    = ""
	)
	date := time.Now()

	// when:
	actualWorkoutDay, err := customer.NewWorkoutDay(customerUUID, customerName, groupUUID, trainerUUID, date)

	// then:
	assertions.Equal(customer.ErrEmptyGroupUUID, err)
	assertions.Empty(actualWorkoutDay)
}

func TestUnmarshalFromDatabaseShouldParseDataWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		workoutDayUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
		trainerUUID    = "d1f9cf76-0acf-445c-b9ee-d9515e20d22d"
		groupUUID      = "aa046e54-1c50-4b85-8a52-764d34c766ef"
		customerUUID   = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
		customerName   = "John Doe"
	)
	date := time.Now()

	// when:
	actualWorkoutDay, err := customer.UnmarshalFromDatabase(customer.UnmarshalFromDatabaseArgs{
		WorkoutDayUUID: workoutDayUUID,
		TrainerUUID:    trainerUUID,
		GroupUUID:      groupUUID,
		CustomerUUID:   customerUUID,
		CustomerName:   customerName,
		Date:           date,
	})

	// then:
	assertions.Nil(err)
	assertions.Equal(groupUUID, actualWorkoutDay.GroupUUID())
	assertions.Equal(trainerUUID, actualWorkoutDay.TrainerUUID())
	assertions.Equal(customerUUID, actualWorkoutDay.CustomerUUID())
	assertions.Equal(customerName, actualWorkoutDay.CustomerName())
	assertions.Equal(date, actualWorkoutDay.Date())
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyTrainerUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		workoutDayUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
		trainerUUID    = ""
		groupUUID      = "aa046e54-1c50-4b85-8a52-764d34c766ef"
		customerUUID   = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
		customerName   = "John Doe"
	)
	date := time.Now()

	// when:
	actualWorkoutDay, err := customer.UnmarshalFromDatabase(customer.UnmarshalFromDatabaseArgs{
		WorkoutDayUUID: workoutDayUUID,
		TrainerUUID:    trainerUUID,
		GroupUUID:      groupUUID,
		CustomerUUID:   customerUUID,
		CustomerName:   customerName,
		Date:           date,
	})

	// then:
	assertions.Equal(err, customer.ErrEmptyTrainerUUID)
	assertions.Empty(actualWorkoutDay)
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyCustomerName_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		workoutDayUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
		trainerUUID    = "d1f9cf76-0acf-445c-b9ee-d9515e20d22d"
		groupUUID      = "aa046e54-1c50-4b85-8a52-764d34c766ef"
		customerUUID   = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
		customerName   = ""
	)
	date := time.Now()

	// when:
	actualWorkoutDay, err := customer.UnmarshalFromDatabase(customer.UnmarshalFromDatabaseArgs{
		WorkoutDayUUID: workoutDayUUID,
		TrainerUUID:    trainerUUID,
		GroupUUID:      groupUUID,
		CustomerUUID:   customerUUID,
		CustomerName:   customerName,
		Date:           date,
	})

	// then:
	assertions.Equal(err, customer.ErrEmptyCustomerName)
	assertions.Empty(actualWorkoutDay)
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyCustomerWorkoutDayUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		workoutDayUUID = ""
		trainerUUID    = "d1f9cf76-0acf-445c-b9ee-d9515e20d22d"
		groupUUID      = "aa046e54-1c50-4b85-8a52-764d34c766ef"
		customerUUID   = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
		customerName   = "John Doe"
	)
	date := time.Now()

	// when:
	SUT, err := customer.UnmarshalFromDatabase(customer.UnmarshalFromDatabaseArgs{
		WorkoutDayUUID: workoutDayUUID,
		TrainerUUID:    trainerUUID,
		GroupUUID:      groupUUID,
		CustomerUUID:   customerUUID,
		CustomerName:   customerName,
		Date:           date,
	})

	// then:
	assertions.Equal(customer.ErrEmptyWorkoutDayUUID, err)
	assertions.Empty(SUT)
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyCustomerUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		workoutDayUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
		customerUUID   = ""
		customerName   = "John Doe"
		trainerUUID    = "aa046e54-1c50-4b85-8a52-764d34c766ef"
		groupUUID      = "aa046e54-1c50-4b85-8a52-764d34c766ef"
	)

	date := time.Now()

	// when:
	workout, err := customer.UnmarshalFromDatabase(customer.UnmarshalFromDatabaseArgs{
		WorkoutDayUUID: workoutDayUUID,
		TrainerUUID:    trainerUUID,
		GroupUUID:      groupUUID,
		CustomerUUID:   customerUUID,
		CustomerName:   customerName,
		Date:           date,
	})

	// then:
	assertions.Equal(customer.ErrEmptyCustomerUUID, err)
	assertions.Empty(workout)
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyGroupUUID_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		workoutDayUUID = "aa046e54-1c50-4b85-8a52-764d34c766ef"
		groupUUID      = ""
		trainerUUID    = "e740726b-6817-4755-a1bb-097c024a76df"
		customerName   = "John Doe"
		customerUUID   = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
	)
	date := time.Now()

	// when:
	SUT, err := customer.UnmarshalFromDatabase(customer.UnmarshalFromDatabaseArgs{
		WorkoutDayUUID: workoutDayUUID,
		TrainerUUID:    trainerUUID,
		GroupUUID:      groupUUID,
		CustomerUUID:   customerUUID,
		CustomerName:   customerName,
		Date:           date,
	})

	// then:
	assertions.Equal(customer.ErrEmptyGroupUUID, err)
	assertions.Empty(SUT)
}

func TestUnmarshalFromDatabaseShouldReturnErrorForEmptyGroupDate_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		WorkoutDayUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
		groupUUID      = "aa046e54-1c50-4b85-8a52-764d34c766ef"
		customerUUID   = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
		trainerUUID    = "02059e83-27f4-46f4-871e-c14b4b9235f0"
		customerName   = "John Doe"
	)
	date := time.Time{}

	// when:
	SUT, err := customer.UnmarshalFromDatabase(customer.UnmarshalFromDatabaseArgs{
		WorkoutDayUUID: WorkoutDayUUID,
		TrainerUUID:    trainerUUID,
		GroupUUID:      groupUUID,
		CustomerUUID:   customerUUID,
		CustomerName:   customerName,
		Date:           date,
	})

	// then:
	assertions.Equal(customer.ErrEmptyGroupDate, err)
	assertions.Empty(SUT)
}
