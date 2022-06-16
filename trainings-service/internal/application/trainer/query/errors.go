package query

import "errors"

var (
	ErrWorkoutGroupNotOwner = errors.New("workout group not owner")
	ErrRepositoryFailure    = errors.New("repository failure")
)
