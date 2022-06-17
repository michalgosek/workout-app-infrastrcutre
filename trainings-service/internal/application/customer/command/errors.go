package command

import "errors"

var (
	ErrRepositoryFailure = errors.New("repository failure")
	ErrResourceNotFound  = errors.New("resource not found")
)
