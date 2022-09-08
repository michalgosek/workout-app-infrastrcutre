package main

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb"
	dbcmd "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	dbqry "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/notifications"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/server"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/ports"
	"log"
	"net/http"
	"time"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	mongoCLI, err := mongodb.NewClient("mongodb://mongo:27017", 5*time.Second)
	if err != nil {
		return err
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()

		err := mongoCLI.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}()

	cmdConfig := dbcmd.Config{
		Database:       "trainings_service_db",
		Collection:     "trainings",
		CommandTimeout: 10 * time.Second,
	}

	insertTrainingGroup := dbcmd.NewInsertTrainingGroupHandler(mongoCLI, cmdConfig)
	updateTrainingGroup := dbcmd.NewUpdateTrainingGroupHandler(mongoCLI, cmdConfig)
	deleteTrainingGroups := dbcmd.NewDeleteTrainingGroupsHandler(mongoCLI, cmdConfig)
	deleteTrainingGroup := dbcmd.NewDeleteTrainingGroupHandler(mongoCLI, cmdConfig)
	queryCfg := dbqry.Config{
		Database:     "trainings_service_db",
		Collection:   "trainings",
		QueryTimeout: 10 * time.Second,
	}

	duplicate := dbqry.NewDuplicateTrainingGroupHandler(mongoCLI, queryCfg)
	trainingGroup := dbqry.NewTrainingGroupHandler(mongoCLI, queryCfg)
	trainerGroup := dbqry.NewTrainerGroupHandler(mongoCLI, queryCfg)
	trainerGroups := dbqry.NewTrainerGroupsHandler(mongoCLI, queryCfg)
	trainerParticipants := dbqry.NewTrainerParticipantsHandler(mongoCLI, queryCfg)
	participantGroups := dbqry.NewParticipantGroupsHandler(mongoCLI, queryCfg)
	allTrainings := dbqry.NewAllTrainingsHandler(mongoCLI, queryCfg)
	notification := notifications.NewService(http.DefaultClient)

	serverCfg := server.DefaultHTTPConfig(":8070", "trainings-service")
	HTTP := ports.NewTrainingsHTTP(&application.Application{
		Commands: application.Commands{
			PlanTrainingGroup:    command.NewPlanTrainingGroupHandler(insertTrainingGroup, duplicate),
			CancelTrainingGroup:  command.NewCancelTrainingGroupHandler(trainingGroup, deleteTrainingGroup, notification),
			CancelTrainingGroups: command.NewCancelTrainingGroupsHandler(deleteTrainingGroups, trainerParticipants, notification),
			UnassignParticipant:  command.NewUnassignParticipantHandler(updateTrainingGroup, trainingGroup),
			AssignParticipant:    command.NewAssignParticipantHandler(updateTrainingGroup, trainingGroup),
			UpdateTrainingGroup:  command.NewUpdateTrainingGroupHandler(updateTrainingGroup, trainingGroup),
		},
		Queries: application.Queries{
			TrainerGroup:      query.NewTrainerGroupHandler(trainerGroup),
			TrainerGroups:     query.NewTrainerGroupsHandler(trainerGroups),
			AllTrainingGroups: query.NewAllTrainingGroupsHandler(allTrainings),
			ParticipantGroups: query.NewParticipantGroupsHandler(participantGroups),
		},
	}, serverCfg.Addr)

	API := HTTP.NewAPI()
	srv := server.NewHTTP(API, serverCfg)
	srv.StartHTTPServer()
	return nil
}
