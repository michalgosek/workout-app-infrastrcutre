package domain

import "time"

type Training struct {
	uuid     string
	userUUID string
	userName string
	data     time.Duration
	canceled bool
}

func NewTraining(uuid string, userUUID string, data time.Duration) (*Training, error) {
	// verification logci

	t := Training{
		uuid:     uuid,
		userUUID: userUUID,
		userName: "",
		data:     data,
		canceled: false,
	}
	return &t, nil
}
