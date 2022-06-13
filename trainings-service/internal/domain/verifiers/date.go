package verifiers

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
	return ErrDateValueViolation
}

var ErrDateValueViolation = errors.New("specified date must be at least 3h before current hour")
