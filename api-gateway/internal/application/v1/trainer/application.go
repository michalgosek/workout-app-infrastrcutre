package trainer

import (
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/v1/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/v1/trainer/query"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Commands struct {
	*command.PlanTrainingHandler
}

type Queries struct {
	*query.TrainingHandler
}

type Application struct {
	Commands
	Queries
}

func NewApplication(HTTP HTTPClient) *Application {
	a := Application{
		Commands: Commands{
			PlanTrainingHandler: command.NewPlanTrainingHandler(HTTP),
		},
		Queries: Queries{
			TrainingHandler: query.NewTrainingHandler(HTTP),
		},
	}
	return &a
}
