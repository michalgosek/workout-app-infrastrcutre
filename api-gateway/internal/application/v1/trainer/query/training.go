package query

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Participant struct {
	Name string
	UUID string
}

type TrainingResponse struct {
	UUID         string
	Name         string
	Description  string
	Date         time.Time
	Limit        int
	Participants []Participant
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type TrainingHandler struct {
	url string
	cli HTTPClient
}

type Training struct {
	UUID        string
	TrainerUUID string
}

func (t *TrainingHandler) Do(ctx context.Context, tr Training) (TrainingResponse, error) {
	url := fmt.Sprintf("http://localhost:8070/api/v1/trainers/%s/trainings/%s", tr.TrainerUUID, tr.UUID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return TrainingResponse{}, err
	}
	res, err := t.cli.Do(req)
	if err != nil {
		return TrainingResponse{}, err
	}
	defer res.Body.Close()

	dec := json.NewDecoder(res.Body)

	var dst TrainingResponse
	err = dec.Decode(&dst)
	if err != nil {
		return TrainingResponse{}, err
	}
	return dst, nil
}

func NewTrainingHandler(c HTTPClient) *TrainingHandler {
	if c == nil {
		panic("nil http client implementation")
	}
	h := TrainingHandler{
		cli: c,
	}
	return &h
}
