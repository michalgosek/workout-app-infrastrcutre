package main

import (
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/notifications"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/server"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/ports/http"
	"log"
	"time"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	repository, err := mongodb.NewRepository(mongodb.Config{
		Addr:       "mongodb://localhost:27017",
		Database:   "trainings_service_db",
		Collection: "trainings",
		Timeouts: mongodb.Timeouts{
			CommandTimeout:    10 * time.Second,
			QueryTimeout:      10 * time.Second,
			ConnectionTimeout: 10 * time.Second,
		},
	})
	if err != nil {
		return err
	}
	defer func() {
		err := repository.Disconnect()
		if err != nil {
			panic(err)
		}
	}()

	notificationService := notifications.NewService()

	serverCfg := server.DefaultHTTPConfig("localhost:8070", "trainings-service")
	HTTP := http.NewTrainingsHTTP(&application.Application{
		Commands: application.Commands{
			PlanTrainingGroup:    command.NewPlanTrainingGroupHandler(repository),
			CancelTrainingGroup:  command.NewCancelTrainingGroupHandler(repository, notificationService),
			CancelTrainingGroups: command.NewCancelTrainingGroupsHandler(repository),
			UnassignParticipant:  command.NewUnassignParticipantHandler(repository),
			AssignParticipant:    command.NewAssignParticipantHandler(repository),
			UpdateTrainingGroup:  command.NewUpdateTrainingGroupHandler(repository),
		},
		Queries: application.Queries{
			TrainerGroup:      query.NewTrainerGroupHandler(repository),
			TrainerGroups:     query.NewTrainerGroupsHandler(repository),
			AllTrainingGroups: query.NewAllTrainingGroupsHandler(repository),
			ParticipantGroups: query.NewParticipantGroupsHandler(repository),
		},
	}, serverCfg.Addr)

	API := HTTP.NewAPI()
	srv := server.NewHTTP(API, serverCfg)
	srv.StartHTTPServer()
	return nil
}
