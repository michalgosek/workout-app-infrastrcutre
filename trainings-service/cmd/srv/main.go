package main

import (
	"fmt"
	cservice "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/customer"
	tservice "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainer"
	trservice "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainings"

	"log"
	"time"

	"github.com/go-chi/chi"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	customcmd "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command"
	trainercmd "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	trainerqry "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/ports/http"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	customerRepository, err := customer.NewCustomerRepository(customer.RepositoryConfig{
		Addr:               "mongodb://localhost:27017",
		Database:           "trainings_service_test",
		CustomerCollection: "customer_schedules",
		CommandTimeout:     10 * time.Second,
		QueryTimeout:       10 * time.Second,
		ConnectionTimeout:  10 * time.Second,
		Format:             "02/01/2006 15:04",
	})
	if err != nil {
		return fmt.Errorf("creating customer repository failed: %v", err)
	}
	customerService, err := cservice.NewCustomerService(customerRepository)
	if err != nil {
		return fmt.Errorf("creating customer service failed: %v", err)
	}

	trainerRepository, err := trainer.NewTrainerRepository(trainer.RepositoryConfig{
		Addr:              "mongodb://localhost:27017",
		Database:          "trainings_service_test",
		TrainerCollection: "trainer_schedules",
		CommandTimeout:    10 * time.Second,
		QueryTimeout:      10 * time.Second,
		ConnectionTimeout: 10 * time.Second,
		Format:            "02/01/2006 15:04",
	})
	if err != nil {
		return fmt.Errorf("creating trainer repository failed: %v", err)
	}
	trainerService, err := tservice.NewTrainerService(trainerRepository)
	if err != nil {
		return fmt.Errorf("creating trainer service failed: %v", err)
	}
	trainingsService, err := trservice.NewService(customerService, trainerService)
	if err != nil {
		return fmt.Errorf("creating trainings service failed: %v", err)
	}

	customerScheduleWorkoutHandler, err := customcmd.NewScheduleWorkoutHandler(trainingsService)
	if err != nil {
		return fmt.Errorf("creating customer schedule workout handler failed: %v", err)
	}
	deleteTrainerWorkoutHandler, err := trainercmd.NewCancelWorkoutHandler(trainingsService)
	if err != nil {
		return fmt.Errorf("creating cancel workout workout handler failed: %v", err)

	}

	app := application.Application{
		Commands: application.Commands{
			CreateTrainerWorkout:    trainercmd.NewScheduleWorkoutHandler(trainerRepository),
			DeleteTrainerWorkout:    deleteTrainerWorkoutHandler,
			DeleteTrainerWorkouts:   trainercmd.NewCancelWorkoutsHandler(trainerRepository),
			UnassignCustomer:        trainercmd.NewUnassignCustomerHandler(customerRepository, trainerRepository),
			CustomerScheduleWorkout: customerScheduleWorkoutHandler,
		},
		Queries: application.Queries{
			GetTrainerWorkout:  trainerqry.NewWorkoutGroupHandler(trainerRepository),
			GetTrainerWorkouts: trainerqry.NewWorkoutGroupsHandler(trainerRepository),
		},
	}
	HTTP := http.NewTrainerWorkoutGroupsHTTP(&app, "02/01/2006 15:04")
	customerHTTP := http.NewCustomerHTTP(&app, "02/01/2006 15:04")

	API := rest.NewRouter()
	API.Route("/api/v1", func(r chi.Router) {
		r.Route("/trainers", func(r chi.Router) {
			r.Route("/{trainerUUID}", func(r chi.Router) {
				r.Route("/workouts", func(r chi.Router) {
					r.Get("/", HTTP.GetTrainerWorkoutGroups())
					r.Post("/", HTTP.CreateTrainerWorkoutGroup())
					r.Delete("/", HTTP.DeleteWorkoutGroups())
					r.Route("/{groupUUID}", func(r chi.Router) {
						r.Get("/", HTTP.GetTrainerWorkoutGroup())
						r.Delete("/", HTTP.DeleteWorkoutGroup())
						r.Route("/customers", func(r chi.Router) {
							//r.Post("/", HTTP.AssignCustomer())
							r.Delete("/{customerUUID}", HTTP.UnassignCustomer())
						})
					})
				})
			})
		})
		r.Route("/customers", func(r chi.Router) {
			r.Route("/{customerUUID}", func(r chi.Router) {
				r.Route("/workouts", func(r chi.Router) {
					r.Post("/", customerHTTP.CreateCustomerWorkout())
				})
			})
		})
	})

	serverCfg := server.DefaultHTTPConfig("localhost:8070", "trainings-service")
	srv := server.NewHTTP(API, serverCfg)
	srv.StartHTTPServer()
	return nil
}
