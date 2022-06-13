package application

import "errors"

var (
	ErrScheduleNotOwner  = errors.New("schedule not owner")
	ErrRepositoryFailure = errors.New("repository failure")
)
