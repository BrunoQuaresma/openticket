package sdk

import (
	"fmt"
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

func (c *Client) DeleteTicket(ticketId int32) (*http.Response, error) {
	httpRes, err := c.Delete("/tickets/" + fmt.Sprint(ticketId))
	return httpRes, err
}

func (c *Client) PatchTicket(ticketId int32, req api.PatchTicketRequest, res *api.PatchTicketResponse) (*http.Response, error) {
	httpRes, err := c.Patch("/tickets/"+fmt.Sprint(ticketId), req, res)
	return httpRes, err
}

func (c *Client) Ticket(ticketId int32, res *api.TicketResponse) (*http.Response, error) {
	httpRes, err := c.Get("/tickets/"+fmt.Sprint(ticketId), res)
	return httpRes, err
}
