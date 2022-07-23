package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/adapters/mongodb"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/application/command"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/application/query"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/ports/http"
	"log"
	"time"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	API := rest.NewRouter()
	repository, err := mongodb.NewRepository(mongodb.Config{
		Addr:       "mongodb://localhost:27017",
		Database:   "users_service",
		Collection: "users",
		Timeouts: mongodb.Timeouts{
			CommandTimeout:    10 * time.Second,
			QueryTimeout:      10 * time.Second,
			ConnectionTimeout: 10 * time.Second,
		},
	})
	if err != nil {
		return fmt.Errorf("users repository creation failed: %s", err)
	}
	defer func() {
		err := repository.Disconnect()
		if err != nil {
			panic(err)
		}
	}()

	app := application.Application{
		Commands: application.Commands{
			RegisterUser: command.NewRegisterHandlerRepository(repository),
		},
		Queries: application.Queries{
			User: query.NewUserHandlerRepository(repository),
		},
	}

	serverCfg := server.DefaultHTTPConfig("localhost:8060", "users-service")
	HTTP := http.NewHTTP(&app, serverCfg.Addr)

	API.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", HTTP.CreateUser())
			r.Get("/{UUID}", HTTP.GetUser())
		})
	})

	srv := server.NewHTTP(API, serverCfg)
	srv.StartHTTPServer()
	return nil
}
