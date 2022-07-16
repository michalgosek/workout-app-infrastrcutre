package trainings

import "errors"

type Participant struct {
	uuid string
	name string
}

func (p Participant) UUID() string {
	return p.uuid
}

func (p Participant) Name() string {
	return p.name
}

func NewParticipant(uuid string, name string) (Participant, error) {
	if uuid == "" {
		return Participant{}, ErrEmptyParticipantUUID
	}
	if name == "" {
		return Participant{}, ErrEmptyParticipantName
	}
	p := Participant{
		uuid: uuid,
		name: name,
	}
	return p, nil
}

var (
	ErrEmptyParticipantUUID = errors.New("specified empty participant uuid")
	ErrEmptyParticipantName = errors.New("specified empty participant name")
)
