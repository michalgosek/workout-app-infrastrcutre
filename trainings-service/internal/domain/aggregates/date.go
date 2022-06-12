package aggregates

import (
	"errors"
	"fmt"
	"time"
)

type WorkoutDate struct {
	hours int
}

func NewWorkoutDate(hours int) WorkoutDate {
	return WorkoutDate{
		hours: hours,
	}
}

func (w *WorkoutDate) Check(providedDate time.Time) error {
	fmt.Println("provied: ", providedDate.String())
	threshold := time.Now().Add(time.Duration(w.hours) * time.Hour)
	fmt.Println("current 3h", threshold.String())

	if providedDate.Equal(threshold) || providedDate.After(threshold) {
		return nil
	}
	return ErrDateAggregateViolation
}

var ErrDateAggregateViolation = errors.New("specified date year should be equal current year")
