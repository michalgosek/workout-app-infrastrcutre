package trainings

import (
	"errors"
	"time"
)

type WorkoutGroup struct {
	uuid         string
	name         string
	description  string
	date         time.Time
	trainer      Trainer
	limit        int
	participants []Participant
}

func (w *WorkoutGroup) UUID() string {
	return w.uuid
}

func (w *WorkoutGroup) Name() string {
	return w.name
}

func (w *WorkoutGroup) Description() string {
	return w.description
}

func (w *WorkoutGroup) Date() time.Time {
	return w.date
}

func (w *WorkoutGroup) Trainer() Trainer {
	return w.trainer
}

func (w *WorkoutGroup) Limit() int {
	return w.limit
}

func (w *WorkoutGroup) Participants() []Participant {
	return w.participants
}

func (w *WorkoutGroup) AssignParticipant(p Participant) error {
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

func (w *WorkoutGroup) UnassignParticipant(UUID string) error {
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
	w.limit = len(filtered)
	return nil
}

func NewWorkoutGroup(uuid, name, desc string, date time.Time, t Trainer) (*WorkoutGroup, error) {
	if uuid == "" {
		return nil, ErrEmptyWorkoutGroupUUID
	}
	if name == "" {
		return nil, ErrEmptyWorkoutGroupName
	}
	if desc == "" {
		return nil, ErrEmptyWorkoutGroupDescription
	}
	g := WorkoutGroup{
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
	ErrEmptyWorkoutGroupUUID        = errors.New("empty workout group uuid")
	ErrEmptyWorkoutGroupName        = errors.New("empty workout group name")
	ErrEmptyWorkoutGroupDescription = errors.New("empty workout group description")
)

type DatabaseWorkoutGroup struct {
	UUID         string
	Name         string
	Description  string
	Limit        int
	Date         time.Time
	Trainer      DatabaseWorkoutGroupTrainer
	Participants []DatabaseWorkoutGroupParticipant
}

type DatabaseWorkoutGroupTrainer struct {
	UUID string
	Name string
}

type DatabaseWorkoutGroupParticipant struct {
	UUID string
	Name string
}

func UnmarshalWorkoutGroupFromDatabase(d DatabaseWorkoutGroup) WorkoutGroup {
	var pp []Participant
	for _, p := range d.Participants {
		pp = append(pp, Participant{
			uuid: p.UUID,
			name: p.Name,
		})
	}
	return WorkoutGroup{
		uuid:        d.UUID,
		name:        d.Name,
		description: d.Description,
		date:        d.Date,
		trainer: Trainer{
			uuid: d.Trainer.UUID,
			name: d.Trainer.Name,
		},
		limit:        d.Limit,
		participants: pp,
	}
}
