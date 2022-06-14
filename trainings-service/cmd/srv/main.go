package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-chi/chi"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/ports"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	cfg := mongodb.Config{
		Addr:              "mongodb://localhost:27017",
		Database:          "trainings_service_test",
		TrainerCollection: "trainer_schedules",
		CommandTimeout:    10 * time.Second,
		QueryTimeout:      10 * time.Second,
		ConnectionTimeout: 10 * time.Second,
		Format:            "02/01/2006 15:04",
	}
	repository, err := mongodb.NewMongoDB(cfg)
	if err != nil {
		return fmt.Errorf("creating repository failed: %v", err)
	}

	service := application.NewTrainerService(repository)
	HTTP := ports.NewHTTP(service, cfg.Format)

	API := rest.NewRouter()
	API.Route("/api/v1/", func(r chi.Router) {
		r.Route("/trainer", func(r chi.Router) {
			r.Post("/group", HTTP.CreateTrainerWorkoutGroup)
			r.Get("/groups", HTTP.GetTrainerWorkoutGroups)
			r.Get("/group", HTTP.GetTrainerWorkoutGroup)
			r.Delete("/group", HTTP.DeleteWorkoutGroup)
			r.Delete("/groups", HTTP.DeleteWorkoutGroups)
		})
	})

	serverCfg := server.DefaultHTTPConfig("localhost:8070", "trainings-service")
	srv := server.NewHTTP(API, serverCfg)
	srv.StartHTTPServer()
	return nil
}
