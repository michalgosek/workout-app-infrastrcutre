package command

import "errors"

var (
	ErrRepositoryFailure      = errors.New("repository failure")
	ErrWorkoutGroupDuplicated = errors.New("duplicate workout group")
	ErrResourceNotFound       = errors.New("resource not found")
)
