package sdk

import (
	"fmt"
	"net/http"

	"github.com/BrunoQuaresma/openticket/api"
)

func (c *Client) CreateUser(req api.CreateUserRequest, res *api.CreateUserResponse) (*http.Response, error) {
	httpRes, err := c.Post("/users", req, res)
	return httpRes, err
}

func (c *Client) DeleteUser(id int32) (*http.Response, error) {
	httpRes, err := c.Delete("/users/" + fmt.Sprint(id))
	return httpRes, err
}
