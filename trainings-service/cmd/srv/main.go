package main

import (
	"fmt"
	"log"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/ports"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	cfg := adapters.MongoDBConfig{
		Addr:              "mongodb://localhost:27017",
		Database:          "trainings_service_test",
		Collection:        "trainer_schedules",
		CommandTimeout:    10 * time.Second,
		QueryTimeout:      10 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	}
	repository, err := adapters.NewTrainerSchedulesMongoDB(cfg)
	if err != nil {
		return fmt.Errorf("creating repository failed: %v", err)
	}

	service := application.NewTrainerService(repository)
	_ = ports.NewHTTP(service)

	API := rest.NewAPI()
	API.SetEndpoints()

	serverCfg := server.DefaultHTTPConfig("localhost:8070", "trainings-service")
	srv := server.NewHTTP(API, serverCfg)
	srv.StartHTTPServer()

	return nil
}
