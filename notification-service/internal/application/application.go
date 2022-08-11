package application

import (
	"notification-service/internal/application/command"
	"notification-service/internal/application/query"
)

type Command struct {
	*command.ClearNotificationsHandler
	*command.InsertNotificationHandler
}

type Query struct {
	*query.AllNotificationsHandler
}

type Application struct {
	*Query
	*Command
}
