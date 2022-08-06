package main

import (
	"github.com/go-chi/chi"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
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

	serverCfg := server.DefaultHTTPConfig("localhost:8070", "trainings-service")
	HTTP := http.NewTrainingsHTTP(&application.Application{
		Commands: application.Commands{
			PlanTrainingGroup:    command.NewPlanTrainingGroupHandler(repository),
			CancelTrainingGroup:  command.NewCancelTrainingGroupHandler(repository),
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
			r.Get("/", HTTP.GetAllTrainingGroups())
		})
		r.Route("/participants", func(r chi.Router) {
			r.Route("/{participantUUID}", func(r chi.Router) {
				r.Get("/", HTTP.GetParticipantGroups())
			})
		})

		r.Route("/trainers", func(r chi.Router) {
			r.Route("/{trainerUUID}", func(r chi.Router) {
				r.Get("/", HTTP.GetTrainerGroups())
				r.Delete("/", HTTP.DeleteTrainerGroups())
				r.Route("/trainings", func(r chi.Router) {
					r.Delete("/", HTTP.DeleteTrainerGroups())
					r.Route("/{trainingUUID}", func(r chi.Router) {
						r.Put("/", HTTP.UpdateTrainingGroup())
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
