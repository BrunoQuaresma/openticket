package sdk

import (
	"net/http"

	"github.com/BrunoQuaresma/openticket/api"
)

func (c *Client) CreateLabel(req api.CreateLabelRequest, res *api.CreateLabelResponse) (*http.Response, error) {
	httpRes, err := c.post("/labels", req, res)
	return httpRes, err
}
