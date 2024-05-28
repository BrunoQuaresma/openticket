package sdk

import (
	"net/http"

	"github.com/BrunoQuaresma/openticket/api"
)

func (c *Client) CreateUser(req api.CreateUserRequest, res *api.CreateUserResponse) (*http.Response, error) {
	httpRes, err := Post(c.url+"/users", req, res)
	return httpRes, err
}
