package sdk

import "github.com/BrunoQuaresma/openticket/api"

type RequestResult[T any] struct {
	StatusCode int
	api.Response[T]
}

type Client struct {
	url string
}

func New(url string) *Client {
	return &Client{url: url}
}
