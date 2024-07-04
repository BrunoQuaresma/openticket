package sdk

import (
	"net/http"

	"github.com/BrunoQuaresma/openticket/api"
)

func (c *Client) Health(res *api.HealthResponse) (*http.Response, error) {
	return c.get("/health", res)
}
