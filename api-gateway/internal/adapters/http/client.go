package http

import (
	"errors"
	"fmt"
	"net/http"
)

type DefaultClient struct {
	cli *http.Client
}

func (c *DefaultClient) Do(req *http.Request) (*http.Response, error) {
	res, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}

	httpSuccessfulResponse := res.StatusCode >= 200 && res.StatusCode <= 299
	if !httpSuccessfulResponse {
		msg := fmt.Sprintf("request call failure, status code: %d", res.StatusCode)
		return nil, errors.New(msg)
	}
	return res, nil
}

func NewDefaultClient() *DefaultClient {
	c := DefaultClient{
		cli: http.DefaultClient,
	}
	return &c
}
