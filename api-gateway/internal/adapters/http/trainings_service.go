package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/api/v1/rest/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/api/v1/rest/trainer/query"
	"net/http"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type TrainingsService struct {
	cli Client
}

func (t *TrainingsService) PlanTraining(ctx context.Context, training command.PlanTrainingCommand) error {
	p := PostTrainingGroup{
		User: PostTrainingUser{
			UUID: training.User.UUID,
			Name: training.User.Name,
			Role: training.User.Role,
		},
		GroupName: training.GroupName,
		GroupDesc: training.GroupDesc,
		Date:      training.Date,
	}
	bb, err := json.Marshal(&p)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(bb)
	url := "http://localhost:8070/api/v1/trainings"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return err
	}

	res, err := t.cli.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (t *TrainingsService) TrainingGroup(ctx context.Context, q query.TrainingQuery) (query.TrainingGroup, error) {
	url := fmt.Sprintf("http://localhost:8070/api/v1/trainers/%s/trainings/%s", q.User.UUID, q.TrainingUUID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return query.TrainingGroup{}, err
	}
	req = req.WithContext(ctx)
	res, err := t.cli.Do(req)
	if err != nil {
		return query.TrainingGroup{}, err
	}
	defer res.Body.Close()

	var result query.TrainingGroup
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&result)
	if err != nil {
		return query.TrainingGroup{}, err
	}
	return result, nil
}

func NewTrainingsService(c Client) (*TrainingsService, error) {
	if c == nil {
		return nil, errors.New("nil HTTP client")
	}
	t := TrainingsService{
		cli: c,
	}
	return &t, nil
}
