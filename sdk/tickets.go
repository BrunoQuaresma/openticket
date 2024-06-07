package sdk

import (
	"net/http"

	"github.com/BrunoQuaresma/openticket/api"
)

func (c *Client) CreateTicket(req api.CreateTicketRequest, res *api.CreateTicketResponse) (*http.Response, error) {
	httpRes, err := c.Post("/tickets", req, res)
	return httpRes, err
}
