package main

import (
	"fmt"
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
	cfg := mongodb.Config{
		Addr:       "mongodb://localhost:27017",
		Database:   "trainings_service_db",
		Collection: "customer_schedules",
		Timeouts: mongodb.Timeouts{
			CommandTimeout:    10 * time.Second,
			QueryTimeout:      10 * time.Second,
			ConnectionTimeout: 10 * time.Second,
		},
	}
	repository, err := mongodb.NewRepository(cfg)
	if err != nil {
		return fmt.Errorf("trainings repository creation failed: %s", err)
	}

	trainingsService := service.NewTrainingsService(repository)

	HTTP := http.NewTrainingsHTTP(&application.Application{
		Commands: application.Commands{
			ScheduleTrainerWorkoutGroup:         command.NewScheduleTrainerWorkoutGroupHandler(trainingsService),
			CancelTrainerWorkoutGroup:           command.NewCancelTrainerWorkoutGroupHandler(trainingsService),
			UnassignParticipantFromWorkoutGroup: command.NewUnassignParticipantHandler(trainingsService),
			AssignParticipantToWorkoutGroup:     command.NewAssignParticipantHandler(trainingsService),
		},
		Queries: application.Queries{
			TrainerWorkoutGroup:  query.NewTrainerWorkoutGroupHandler(repository),
			TrainerWorkoutGroups: query.NewTrainerWorkoutGroupsHandler(repository),
		},
	})

	API := rest.NewRouter()
	API.Route("/api/v1", func(r chi.Router) {
		r.Route("/trainers", func(r chi.Router) {
			r.Route("/{trainerUUID}", func(r chi.Router) {
				r.Route("/workouts", func(r chi.Router) {
					r.Get("/", HTTP.GetTrainerWorkoutGroups())
					r.Post("/", HTTP.CreateTrainerWorkoutGroup())
					//r.Delete("/", HTTP.DeleteWorkoutGroups())
					r.Route("/{groupUUID}", func(r chi.Router) {
						r.Get("/", HTTP.GetTrainerWorkoutGroup())
						r.Delete("/", HTTP.DeleteTrainerWorkoutGroup())
						//	r.Route("/customers", func(r chi.Router) {
						r.Post("/", HTTP.AssignParticipantToWorkoutGroup())
						//		r.Delete("/{customerUUID}", HTTP.UnassignCustomer())
						//	})
					})
				})
			})
		})
	})

	serverCfg := server.DefaultHTTPConfig("localhost:8070", "trainings-service")
	srv := server.NewHTTP(API, serverCfg)
	srv.StartHTTPServer()
	return nil
}
