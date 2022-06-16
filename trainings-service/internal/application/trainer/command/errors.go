package command

import "errors"

var (
	ErrWorkoutGroupNotOwner = errors.New("schedule not owner")
	ErrResourceNotFound     = errors.New("resource not found")
	ErrRepositoryFailure    = errors.New("repository failure")
)
