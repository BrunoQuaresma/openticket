package sdk

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/BrunoQuaresma/openticket/api"
)

func (c *Client) CreateTicket(req api.CreateTicketRequest, res *api.CreateTicketResponse) (*http.Response, error) {
	httpRes, err := c.post("/tickets", req, res)
	return httpRes, err
}

func (c *Client) Tickets(res *api.TicketsResponse, urlValues *url.Values) (*http.Response, error) {
	var searchQuery string
	if urlValues != nil {
		searchQuery = "?" + urlValues.Encode()
	}
	httpRes, err := c.get("/tickets"+searchQuery, res)
	return httpRes, err
}

func (c *Client) DeleteTicket(ticketId int32) (*http.Response, error) {
	httpRes, err := c.delete("/tickets/" + fmt.Sprint(ticketId))
	return httpRes, err
}

func (c *Client) PatchTicket(ticketId int32, req api.PatchTicketRequest, res *api.PatchTicketResponse) (*http.Response, error) {
	httpRes, err := c.patch("/tickets/"+fmt.Sprint(ticketId), req, res)
	return httpRes, err
}

func (c *Client) Ticket(ticketId int32, res *api.TicketResponse) (*http.Response, error) {
	httpRes, err := c.get("/tickets/"+fmt.Sprint(ticketId), res)
	return httpRes, err
}

func (c *Client) PatchTicketStatus(ticketId int32, req api.PatchTicketStatusRequest, res *api.PatchTicketStatusResponse) (*http.Response, error) {
	httpRes, err := c.patch("/tickets/"+fmt.Sprint(ticketId)+"/status", req, res)
	return httpRes, err
}
