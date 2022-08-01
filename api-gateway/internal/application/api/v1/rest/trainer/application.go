package trainer

import (
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/api/v1/rest/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/api/v1/rest/trainer/query"
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
