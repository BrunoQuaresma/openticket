package sdk

import (
	"net/http"

	"github.com/BrunoQuaresma/openticket/api"
)

func (c *Client) Status(res *api.StatusResponse) (*http.Response, error) {
	return c.get("/status", res)
}
