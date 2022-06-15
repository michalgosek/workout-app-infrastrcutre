package application

import "errors"

var (
	ErrScheduleNotOwner  = errors.New("schedule not owner")
	ErrEmptyTrainerUUID  = errors.New("empty trainerUUID")
	ErrRepositoryFailure = errors.New("repository failure")
)
