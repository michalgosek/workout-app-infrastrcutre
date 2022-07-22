package command

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type PlanTraining struct {
	TrainerUUID string `json:"trainer_uuid"`
	TrainerName string `json:"trainer_name"`
	GroupName   string `json:"group_name"`
	GroupDesc   string `json:"group_desc"`
	Date        string `json:"date"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type PlanTrainingHandler struct {
	url string
	cli HTTPClient
}

func (p *PlanTrainingHandler) Do(ctx context.Context, t PlanTraining) error {
	bb, err := json.Marshal(t)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(bb)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.url, body)
	if err != nil {
		return err
	}

	res, err := p.cli.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	httpSuccessfulResponse := res.StatusCode >= 200 && res.StatusCode <= 299
	if !httpSuccessfulResponse {
		msg := fmt.Sprintf("request call failure, status code: %d", res.StatusCode)
		return errors.New(msg)
	}
	return nil
}

func NewPlanTrainingHandler(c HTTPClient) *PlanTrainingHandler {
	if c == nil {
		panic("nil http client implementation")
	}
	h := PlanTrainingHandler{
		url: "http://localhost:8070/api/v1/trainings",
		cli: c,
	}
	return &h
}
