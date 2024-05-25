package sdk

import "github.com/BrunoQuaresma/openticket/api"

type RequestResult struct {
	StatusCode int
	api.Response
}

type Client struct {
	url string
}

func New(url string) *Client {
	return &Client{url: url}
}
