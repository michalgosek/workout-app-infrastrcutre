package domain

import "time"

type DatabaseUser struct {
	UUID           string
	Active         bool
	Role           string
	Name           string
	LastActiveDate time.Time
}

func UnmarshalUserFromDatabase(d DatabaseUser) User {
	return User{
		uuid:           d.UUID,
		name:           d.Name,
		lastActiveDate: d.LastActiveDate,
		role:           d.Role,
		active:         d.Active,
	}
}
