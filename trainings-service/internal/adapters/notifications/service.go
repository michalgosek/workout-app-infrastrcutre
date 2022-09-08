package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/authorization"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/command"
	"net/http"
	"time"
)

const (
	NotificationServiceEndpoint = "http://notifications-service:8060/api/v1/notifications"
)

//go:generate mockery --name=HTTPClient --case underscore --with-expecter
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type config struct {
	commandTimeout time.Duration
}

type Service struct {
	cli HTTPClient
	cfg config
}

func (s *Service) CreateNotification(ctx context.Context, n command.Notification) error {
	bb, err := json.Marshal(n)
	if err != nil {
		return err
	}
	buff := bytes.NewBuffer(bb)

	ctx, cancel := context.WithTimeout(ctx, s.cfg.commandTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/%s", NotificationServiceEndpoint, n.UserUUID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buff)
	if err != nil {
		return err
	}

	v := ctx.Value(authorization.Token)
	if v == nil {
		return errors.New("missing bearer token")
	}
	token, ok := v.(string)
	if !ok {
		return errors.New("bearer token type invalid")
	}
	req.Header.Add(authorization.AuthorizationHeaderKey, token)

	res, err := s.cli.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusCreated {
		return errors.New(http.StatusText(res.StatusCode))
	}
	return nil
}

func NewService(c HTTPClient) *Service {
	s := Service{
		cli: c,
		cfg: config{commandTimeout: 10 * time.Second},
	}
	return &s
}
