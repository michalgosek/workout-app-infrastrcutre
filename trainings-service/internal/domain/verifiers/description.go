package verifiers

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
		return ErrScheduleDescriptionExceeded
	}
	return nil
}

var ErrScheduleDescriptionExceeded = errors.New("description limit exceeded")
