package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/v1/users/query"
	"net/http"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type UsersService struct {
	cli Client
}

func (u *UsersService) User(ctx context.Context, UUID string) (query.User, error) {
	url := fmt.Sprintf("http://localhost:8060/api/v1/users/%s", UUID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	res, err := u.cli.Do(req)
	if err != nil {
		return query.User{}, nil
	}
	defer res.Body.Close()

	var user query.User
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&user)
	if err != nil {
		return query.User{}, err
	}
	return user, nil
}

func NewUsersService(c Client) (*UsersService, error) {
	if c == nil {
		return nil, errors.New("nil HTTP client")
	}
	u := UsersService{
		cli: c,
	}
	return &u, nil
}
