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
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query"
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

	app := application.Application{
		Commands: application.Commands{
			CreateTrainerWorkout:  command.NewCreateWorkoutHandler(repository),
			DeleteTrainerWorkout:  command.NewWorkoutDeleteHandler(repository),
			DeleteTrainerWorkouts: command.NewWorkoutsDeleteHandler(repository),
			AssignCustomer:        command.NewAssignCustomerHandler(repository),
		},
		Queries: application.Queries{
			GetTrainerWorkout:  query.NewGetWorkoutHandler(repository),
			GetTrainerWorkouts: query.NewGetWorkoutsHandler(repository),
		},
	}
	HTTP := http.NewTrainerWorkoutGroupsHTTP(&app, cfg.Format)
	API := rest.NewRouter()

	API.Route("/api/v1", func(r chi.Router) {
		r.Route("/trainers", func(r chi.Router) {
			r.Route("/{trainerUUID}", func(r chi.Router) {
				r.Route("/workouts", func(r chi.Router) {
					r.Get("/", HTTP.GetTrainerWorkoutGroups())
					r.Get("/{workoutUUID}", HTTP.GetTrainerWorkoutGroup())
					r.Post("/", HTTP.CreateTrainerWorkoutGroup())
					r.Delete("/", HTTP.DeleteWorkoutGroups())
					r.Delete("/{workoutUUID}", HTTP.DeleteWorkoutGroup())
				})
			})
		})
	})

	//r.Route("/customers", func(r chi.Router) {
	//	r.Route("/{customerId}", func(r chi.Router) {
	//		r.Route("/workouts", func(r chi.Router) {
	//			r.Get("/")
	//		})
	//
	//	})
	//
	//})

	serverCfg := server.DefaultHTTPConfig("localhost:8070", "trainings-service")
	srv := server.NewHTTP(API, serverCfg)
	srv.StartHTTPServer()
	return nil
}
