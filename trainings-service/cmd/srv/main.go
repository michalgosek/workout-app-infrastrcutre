package main

import (
	"github.com/go-chi/chi"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/service"
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
	repository := mongodb.NewRepository(mongodb.Config{
		Addr:       "mongodb://localhost:27017",
		Database:   "trainings_service_db",
		Collection: "trainings",
		Timeouts: mongodb.Timeouts{
			CommandTimeout:    10 * time.Second,
			QueryTimeout:      10 * time.Second,
			ConnectionTimeout: 10 * time.Second,
		},
	})
	defer func() {
		err := repository.Disconnect()
		if err != nil {
			panic(err)
		}
	}()

	trainingsService := service.NewTrainingsService(repository)
	serverCfg := server.DefaultHTTPConfig("localhost:8070", "trainings-service")
	HTTP := http.NewTrainingsHTTP(&application.Application{
		Commands: application.Commands{
			PlanTrainingGroup:    command.NewPlanTrainingGroupHandler(trainingsService),
			CancelTrainingGroup:  command.NewCancelTrainingGroupHandler(trainingsService),
			CancelTrainingGroups: command.NewCancelTrainingGroupsHandler(trainingsService),
			UnassignParticipant:  command.NewUnassignParticipantHandler(trainingsService),
			AssignParticipant:    command.NewAssignParticipantHandler(trainingsService),
		},
		Queries: application.Queries{
			TrainingGroup:  query.NewTrainingGroupHandler(repository),
			TrainingGroups: query.NewTrainingGroupsHandlerHandler(repository),
		},
	}, serverCfg.Addr)

	API := newAPI(HTTP)
	srv := server.NewHTTP(API, serverCfg)
	srv.StartHTTPServer()
	return nil
}

func newAPI(HTTP *http.Trainings) chi.Router {
	API := rest.NewRouter()
	API.Route("/api/v1", func(r chi.Router) {
		r.Route("/trainings", func(r chi.Router) {
			r.Post("/", HTTP.CreateTrainingGroup())
		})
		r.Route("/trainers", func(r chi.Router) {
			r.Route("/{trainerUUID}", func(r chi.Router) {
				r.Get("/", HTTP.GetTrainerGroups())
				r.Delete("/", HTTP.DeleteTrainerGroups())
				r.Route("/trainings", func(r chi.Router) {
					r.Delete("/", HTTP.DeleteTrainerGroups())
					r.Route("/{trainingUUID}", func(r chi.Router) {
						r.Get("/", HTTP.GetTrainerGroup())
						r.Delete("/", HTTP.DeleteTrainerGroup())
						r.Route("/participants", func(r chi.Router) {
							r.Post("/", HTTP.AssignParticipant())
							r.Route("/{participantUUID}", func(r chi.Router) {
								r.Delete("/", HTTP.UnassignParticipant())
							})
						})
					})
				})
			})
		})
	})
	return API
}
