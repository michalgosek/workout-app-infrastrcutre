package main

import (
	"github.com/go-chi/chi"
	adapters "github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/adapters/http"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/api/v1/rest/trainer"
	trainercmd "github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/api/v1/rest/trainer/command"
	trainerquery "github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/api/v1/rest/trainer/query"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/server"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/server/rest"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/ports"
	"log"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	API := rest.NewRouter()
	serverCfg := server.DefaultHTTPConfig("localhost:8080", "application-gateway")
	HTTPCli := adapters.NewDefaultClient()
	t, err := adapters.NewTrainingsService(HTTPCli)
	if err != nil {
		return err
	}

	trainerApp := trainer.Application{
		Commands: trainer.Commands{
			PlanTrainingHandler: trainercmd.NewPlanTrainingHandler(t),
		},
		Queries: trainer.Queries{
			TrainingHandler: trainerquery.NewTrainingHandler(t),
		},
	}

	err = setTrainerRoutes(API, &trainerApp)
	if err != nil {
		return err
	}
	srv := server.NewHTTP(API, serverCfg)
	srv.StartHTTPServer()
	return nil
}

func setTrainerRoutes(r chi.Router, t *trainer.Application) error {
	HTTP := ports.NewTrainerHTTP(t)

	r.Route("/api/v1/trainings", func(r chi.Router) {
		r.Route("/trainer", func(r chi.Router) {
			r.Post("/", HTTP.CreateTraining())
		})
		r.Route("/{trainingUUID}", func(r chi.Router) {
			r.Route("/trainers", func(r chi.Router) {
				r.Route("/{userUUID}", func(r chi.Router) {
					r.Get("/", HTTP.GetTraining())
				})
			})
		})
	})
	return nil
}
