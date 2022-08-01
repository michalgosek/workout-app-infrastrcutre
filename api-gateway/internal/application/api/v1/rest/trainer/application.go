package trainer

import (
	adapters "github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/adapters/http"
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

func NewApplication(t *adapters.TrainingsService, u *adapters.UsersService) (*Application, error) {
	planTraining, err := command.NewPlanTrainingHandler(u, t)
	if err != nil {
		return nil, err
	}
	training, err := query.NewTrainingHandler(u, t)
	if err != nil {
		return nil, err
	}

	a := Application{
		Commands: Commands{
			PlanTrainingHandler: planTraining,
		},
		Queries: Queries{
			TrainingHandler: training,
		},
	}
	return &a, nil
}
