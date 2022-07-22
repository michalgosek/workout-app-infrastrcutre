package main

import (
	"github.com/go-chi/chi"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/v1/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/ports"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"
	"log"
	"net/http"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	API := rest.NewRouter()
	serverCfg := server.DefaultHTTPConfig("localhost:8080", "application-gateway")

	setTrainerRoutes(API)

	srv := server.NewHTTP(API, serverCfg)
	srv.StartHTTPServer()
	return nil
}

func setTrainerRoutes(r chi.Router) {
	HTTPCli := http.DefaultClient
	API := trainer.NewApplication(HTTPCli)
	HTTP := ports.NewTrainerHTTP(API)

	r.Route("/api/v1/trainings", func(r chi.Router) {
		r.Route("/trainer", func(r chi.Router) {
			r.Post("/", HTTP.CreateTraining())
		})
		r.Route("/{trainingUUID}", func(r chi.Router) {
			r.Route("/trainers", func(r chi.Router) {
				r.Route("/{trainerUUID}", func(r chi.Router) {
					r.Get("/", HTTP.GetTraining())
				})
			})
		})
	})
}
