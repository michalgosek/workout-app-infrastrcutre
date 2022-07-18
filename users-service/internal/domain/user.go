package domain

import (
	"errors"
	"time"
)

type User struct {
	uuid           string
	active         bool
	role           string
	name           string
	lastActiveDate time.Time
}

func (u *User) UUID() string {
	return u.uuid
}

func (u *User) Role() string {
	return u.role
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Active() bool {
	return u.active
}

func (u *User) LastActiveDate() time.Time {
	return u.lastActiveDate
}

func NewUser(uuid, role, name string) (*User, error) {
	if uuid == "" {
		return nil, errors.New("empty user uuid")
	}
	if role == "" {
		return nil, errors.New("empty user role")
	}
	if name == "" {
		return nil, errors.New("empty user name")
	}
	u := User{
		uuid:   uuid,
		active: true,
		role:   role,
		name:   name,
	}
	return &u, nil
}
