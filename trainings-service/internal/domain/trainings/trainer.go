package trainings

import "errors"

type Trainer struct {
	uuid          string
	name          string
	workoutGroups []WorkoutGroup
}

func (t Trainer) UUID() string {
	return t.uuid
}

func (t Trainer) Name() string {
	return t.name
}

func (t Trainer) WorkoutGroups() []WorkoutGroup {
	return t.workoutGroups
}

func NewTrainer(uuid, name string) (Trainer, error) {
	if uuid == "" {
		return Trainer{}, ErrEmptyTrainerUUID
	}
	if name == "" {
		return Trainer{}, ErrEmptyTrainerName
	}
	t := Trainer{
		uuid: uuid,
		name: name,
	}
	return t, nil
}

var (
	ErrEmptyTrainerUUID = errors.New("empty trainer uuid")
	ErrEmptyTrainerName = errors.New("empty trainer name")
)
