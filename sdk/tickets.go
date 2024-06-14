package sdk

import (
	"net/http"
	"net/url"

	"github.com/BrunoQuaresma/openticket/api"
)

func (c *Client) CreateTicket(req api.CreateTicketRequest, res *api.CreateTicketResponse) (*http.Response, error) {
	httpRes, err := c.Post("/tickets", req, res)
	return httpRes, err
}

func (c *Client) Tickets(res *api.TicketsResponse, urlValues *url.Values) (*http.Response, error) {
	httpRes, err := c.Get("/tickets?"+urlValues.Encode(), res)
	return httpRes, err
}
