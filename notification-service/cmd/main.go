package main

import (
	"log"
	"notification-service/internal/adapters/mongodb/client"
	db_cmd "notification-service/internal/adapters/mongodb/command"
	db_query "notification-service/internal/adapters/mongodb/query"
	"notification-service/internal/application"
	"notification-service/internal/application/command"
	"notification-service/internal/application/query"
	"notification-service/internal/application/server"
	"notification-service/internal/ports"

	"time"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	mongoCLI, err := client.New("mongodb://localhost:27017", 5*time.Second)
	if err != nil {
		return err
	}

	const db = "notification_service_db"

	insertionRepository := db_cmd.NewInsertNotificationHandler(mongoCLI, db_cmd.InsertNotificationHandlerConfig{
		Database:       db,
		CommandTimeout: 10 * time.Second,
	})
	cleanupRepository := db_cmd.NewClearNotificationsHandler(mongoCLI, db_cmd.ClearNotificationsHandlerConfig{
		Database:       db,
		CommandTimeout: 10 * time.Second,
	})
	queryAllRepository := db_query.NewAllNotificationHandler(mongoCLI, db_query.AllNotificationHandlerConfig{
		Database:       db,
		CommandTimeout: 10 * time.Second,
	})

	app := application.Application{
		Query: &application.Query{
			AllNotificationsHandler: query.NewAllNotificationsHandler(queryAllRepository),
		},
		Command: &application.Command{
			ClearNotificationsHandler: command.NewClearNotificationsHandler(cleanupRepository),
			InsertNotificationHandler: command.NewInsertNotificationHandler(insertionRepository),
		},
	}

	HTTP, err := ports.NewHTTP(&app)
	if err != nil {
		return err
	}

	API := HTTP.NewAPI()
	serverCfg := server.DefaultHTTPConfig("localhost:8060", "notification-service")
	srv := server.NewHTTP(API, serverCfg)
	srv.StartHTTPServer()

	return nil
}
