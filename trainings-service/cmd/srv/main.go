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
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/ports/http"
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
	trainerWorkoutGroupsHTTP := http.NewTrainerWorkoutGroupsHTTP(service, cfg.Format)

	API := rest.NewRouter()
	API.Route("/api/v1/", func(r chi.Router) {
		r.Route("/trainer", func(r chi.Router) {
			r.Post("/group", trainerWorkoutGroupsHTTP.CreateTrainerWorkoutGroup)
			r.Get("/groups", trainerWorkoutGroupsHTTP.GetTrainerWorkoutGroups)
			r.Get("/group", trainerWorkoutGroupsHTTP.GetTrainerWorkoutGroup)
			r.Delete("/group", trainerWorkoutGroupsHTTP.DeleteWorkoutGroup)
			r.Delete("/groups", trainerWorkoutGroupsHTTP.DeleteWorkoutGroups)
		})
	})

	serverCfg := server.DefaultHTTPConfig("localhost:8070", "trainings-service")
	srv := server.NewHTTP(API, serverCfg)
	srv.StartHTTPServer()
	return nil
}
