package sdk

import (
	"net/http"

	"github.com/BrunoQuaresma/openticket/api"
)

func (c *Client) CreateTicket(req api.CreateTicketRequest, res *api.CreateTicketResponse) (*http.Response, error) {
	httpRes, err := c.Post("/tickets", req, res)
	return httpRes, err
}

func (c *Client) Tickets(res *api.TicketsResponse) (*http.Response, error) {
	httpRes, err := c.Get("/tickets", res)
	return httpRes, err
}
