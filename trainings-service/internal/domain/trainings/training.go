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

func (w *TrainingGroup) UUID() string {
	return w.uuid
}

func (w *TrainingGroup) Name() string {
	return w.name
}

func (w *TrainingGroup) Description() string {
	return w.description
}

func (w *TrainingGroup) Date() time.Time {
	return w.date
}

func (w *TrainingGroup) Trainer() Trainer {
	return w.trainer
}

func (w *TrainingGroup) Limit() int {
	return w.limit
}

func (w *TrainingGroup) IsTrainingDateDuplicated(date time.Time) bool {
	return w.date.Equal(date)
}

func (w *TrainingGroup) Participants() []Participant {
	return w.participants
}

func (w *TrainingGroup) AssignParticipant(p Participant) error {
	if w.limit == 0 {
		return errors.New("workout group participants limit exceeded")
	}
	if len(w.participants) == 0 {
		w.participants = append(w.participants, p)
		w.limit--
		return nil
	}
	for _, c := range w.participants {
		if c.UUID() == p.uuid {
			return errors.New("participant with specified uuid already assigned to workout group")
		}
	}
	w.participants = append(w.participants, p)
	w.limit--
	return nil
}

func (w *TrainingGroup) UnassignParticipant(UUID string) error {
	if len(w.participants) == 0 {
		return errors.New("workout group participants not found")
	}
	var filtered []Participant
	for _, c := range w.participants {
		if c.UUID() != UUID {
			filtered = append(filtered, c)
		}
	}
	w.participants = filtered
	w.limit += 1
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
