package sdk

import (
	"fmt"
	"net/http"

	"github.com/BrunoQuaresma/openticket/api"
)

func (c *Client) CreateUser(req api.CreateUserRequest, res *api.CreateUserResponse) (*http.Response, error) {
	httpRes, err := c.post("/users", req, res)
	return httpRes, err
}

func (c *Client) DeleteUser(id int32) (*http.Response, error) {
	httpRes, err := c.delete("/users/" + fmt.Sprint(id))
	return httpRes, err
}

func (c *Client) PatchUser(id int32, req api.PatchUserRequest, res *api.PatchUserResponse) (*http.Response, error) {
	httpRes, err := c.patch("/users/"+fmt.Sprint(id), req, res)
	return httpRes, err
}
