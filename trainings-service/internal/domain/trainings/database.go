package trainings

import (
	"time"
)

type Error struct {
	resource bool
	message  string
}

func (e Error) Resource() bool {
	return e.resource
}

func (e Error) Error() string {
	return e.message
}

func NewError(msg string, resource bool) Error {
	return Error{
		resource: resource,
		message:  msg,
	}
}

func IsErrResourceNotFound(err error) bool {
	type resource interface {
		Resource() bool
	}
	re, ok := err.(resource)
	return ok && re.Resource()
}

type DatabaseTrainingGroup struct {
	UUID         string
	Name         string
	Description  string
	Limit        int
	Date         time.Time
	Trainer      DatabaseTrainingGroupTrainer
	Participants []DatabaseTrainingGroupParticipant
}

type DatabaseTrainingGroupTrainer struct {
	UUID string
	Name string
}

type DatabaseTrainingGroupParticipant struct {
	UUID string
	Name string
}

func UnmarshalTrainingGroupFromDatabase(d DatabaseTrainingGroup) TrainingGroup {
	var pp []Participant
	for _, p := range d.Participants {
		pp = append(pp, Participant{
			uuid: p.UUID,
			name: p.Name,
		})
	}
	return TrainingGroup{
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
