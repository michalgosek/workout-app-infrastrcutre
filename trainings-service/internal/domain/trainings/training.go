package trainings

import (
	"errors"
	"time"
)

type TrainingGroup struct {
	uuid         string
	name         string
	description  string
	date         time.Time
	trainer      Trainer
	limit        int
	participants []Participant
}

func (t *TrainingGroup) UUID() string {
	return t.uuid
}

func (t *TrainingGroup) Name() string {
	return t.name
}

func (t *TrainingGroup) Description() string {
	return t.description
}

func (t *TrainingGroup) Date() time.Time {
	return t.date
}

func (t *TrainingGroup) Trainer() Trainer {
	return t.trainer
}

func (t *TrainingGroup) Limit() int {
	return t.limit
}

func (t *TrainingGroup) IsOwnedByTrainer(UUID string) bool {
	return UUID == t.trainer.uuid
}

func (t *TrainingGroup) Participants() []Participant {
	return t.participants
}

func (t *TrainingGroup) AssignParticipant(p Participant) error {
	if t.limit == 0 {
		return errors.New("workout group participants limit exceeded")
	}
	if len(t.participants) == 0 {
		t.participants = append(t.participants, p)
		t.limit--
		return nil
	}
	for _, c := range t.participants {
		if c.UUID() == p.uuid {
			return errors.New("participant with specified uuid already assigned to workout group")
		}
	}
	t.participants = append(t.participants, p)
	t.limit--
	return nil
}

func (t *TrainingGroup) UnassignParticipant(UUID string) error {
	if len(t.participants) == 0 {
		return errors.New("workout group participants not found")
	}
	var filtered []Participant
	for _, c := range t.participants {
		if c.UUID() != UUID {
			filtered = append(filtered, c)
		}
	}
	t.participants = filtered
	t.limit += 1
	return nil
}

func NewTrainingGroup(uuid, name, desc string, date time.Time, t Trainer) (*TrainingGroup, error) {
	if uuid == "" {
		return nil, ErrEmptyTrainingGroupUUID
	}
	if name == "" {
		return nil, ErrEmptyTrainingGroupName
	}
	if desc == "" {
		return nil, ErrEmptyTrainingGroupDescription
	}
	g := TrainingGroup{
		uuid:        uuid,
		name:        name,
		description: desc,
		date:        date,
		limit:       10,
		trainer:     t,
	}
	return &g, nil
}

var (
	ErrEmptyTrainingGroupUUID        = errors.New("empty workout group uuid")
	ErrEmptyTrainingGroupName        = errors.New("empty workout group name")
	ErrEmptyTrainingGroupDescription = errors.New("empty workout group description")
)
