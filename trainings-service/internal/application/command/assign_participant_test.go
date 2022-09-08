package command_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/command/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestShouldAssignParticipantWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	updateTrainingGroupRepository := mocks.NewUpdateTrainingGroupRepository(t)
	trainingGroupRepository := mocks.NewTrainingGroupRepository(t)

	// given:
	ctx := context.Background()
	SUT := command.NewAssignParticipantHandler(updateTrainingGroupRepository, trainingGroupRepository)

	trainer, _ := trainings.NewTrainer("05d24119-131d-4362-b4e8-200fe017ea09", "Jerry Smith")
	participant, _ := trainings.NewParticipant("aeca78ac-78fc-40d1-92eb-cb1de4e81d2c", "John Doe")
	training, _ := trainings.NewTrainingGroup("c4d8b177-6666-47a9-87e4-c3f609efc97b", "training", "training desc", time.Now(), trainer)

	trainingWithAssignedParticipant := *training
	_ = trainingWithAssignedParticipant.AssignParticipants(participant)

	trainingGroupRepository.EXPECT().TrainingGroup(ctx, training.UUID()).Return(*training, nil)
	updateTrainingGroupRepository.EXPECT().UpdateTrainingGroup(ctx, &trainingWithAssignedParticipant).Return(nil)

	// when:
	err := SUT.Do(ctx, command.AssignParticipant{
		TrainerUUID:  trainer.UUID(),
		TrainingUUID: training.UUID(),
		Participant:  participant,
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, updateTrainingGroupRepository, trainingGroupRepository)
}
