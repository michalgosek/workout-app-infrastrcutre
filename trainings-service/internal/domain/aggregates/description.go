package aggregates

import "errors"

type WorkoutDescription struct {
	length int
}

func NewWorkoutDescription(length int) WorkoutDescription {
	return WorkoutDescription{
		length: length,
	}
}

func (w *WorkoutDescription) Check(s string) error {
	if len(s) > w.length {
		return ErrWorkoutSessionDescriptionExceeded
	}
	return nil
}

var ErrWorkoutSessionDescriptionExceeded = errors.New("description limit exceeded")
