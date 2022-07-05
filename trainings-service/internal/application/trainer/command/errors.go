package command

import "errors"

var (
	ErrWorkoutGroupNotOwner       = errors.New("schedule not owner")
	ErrWorkoutGroupDateDuplicated = errors.New("duplicate workout group date")
	ErrResourceNotFound           = errors.New("resource not found")
	ErrRepositoryFailure          = errors.New("repository failure")
)
