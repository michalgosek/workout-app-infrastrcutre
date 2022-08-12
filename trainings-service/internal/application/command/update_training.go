package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"time"
)

type UpdateTrainingGroup struct {
	TrainingUUID string
	TrainerUUID  string
	Name         string
	Description  string
	Date         time.Time
}

func (u *UpdateTrainingGroup) Validate() error {
	params := map[string]string{
		"Trainer UUID":               u.TrainerUUID,
		"Training UUID":              u.TrainingUUID,
		"Training Group Name":        u.Name,
		"Training Group Description": u.Description,
	}
	for k, v := range params {
		if v == "" {
			msg := fmt.Sprintf("empty value for %s", k)
			return errors.New(msg)
		}
	}
	if u.Date.IsZero() {
		return errors.New("zero value date has been provided")
	}
	now := time.Now()
	if u.Date.Before(now) {
		return errors.New("an earlier date has been provided")
	}
	return nil
}

type UpdateTrainingGroupRepository interface {
	UpdateTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error
}

type UpdateTrainingGroupHandler struct {
	command UpdateTrainingGroupRepository
	query   TrainingGroupRepository
}

func (u *UpdateTrainingGroupHandler) Do(ctx context.Context, cmd UpdateTrainingGroup) error {
	found, err := u.query.TrainingGroup(ctx, cmd.TrainingUUID)
	if err != nil {
		return err
	}
	if !found.IsOwnedByTrainer(cmd.TrainerUUID) {
		return ErrTrainingNotOwnedByTrainer
	}

	updated, err := trainings.NewTrainingGroup(found.UUID(), cmd.Name, cmd.Description, cmd.Date, found.Trainer())
	if err != nil {
		return err
	}

	err = updated.AssignParticipants(found.Participants()...)
	if err != nil {
		return err
	}
	err = u.command.UpdateTrainingGroup(ctx, updated)
	if err != nil {
		return err
	}
	return nil
}

func NewUpdateTrainingGroupHandler(cmd UpdateTrainingGroupRepository, query TrainingGroupRepository) *UpdateTrainingGroupHandler {
	if cmd == nil {
		panic("nil update training group repository")
	}
	if query == nil {
		panic("nil training group repository")
	}
	h := UpdateTrainingGroupHandler{
		command: cmd,
		query:   query,
	}
	return &h
}

var ErrTrainingNotOwnedByTrainer = errors.New("training not owned by trainer")
