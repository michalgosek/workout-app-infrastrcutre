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
	email          string
	lastActiveDate time.Time
}

func (u *User) IsEmailAlreadyRegistered(email string) bool {
	return u.email == email
}

func (u *User) Email() string {
	return u.email
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

func NewUser(uuid, role, name, email string) (*User, error) {
	if uuid == "" {
		return nil, errors.New("empty user uuid")
	}
	if role == "" {
		return nil, errors.New("empty user role")
	}
	if name == "" {
		return nil, errors.New("empty user name")
	}
	if email == "" {
		return nil, errors.New("empty user email address")
	}
	u := User{
		uuid:   uuid,
		active: true,
		role:   role,
		email:  email,
		name:   name,
	}
	return &u, nil
}
