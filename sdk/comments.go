package sdk

import (
	"fmt"
	"net/http"

	"github.com/BrunoQuaresma/openticket/api"
)

func (c *Client) CreateComment(ticketId int32, req api.CreateCommentRequest, res *api.CreateCommentResponse) (*http.Response, error) {
	httpRes, err := c.Post("/tickets/"+fmt.Sprint(ticketId)+"/comments", req, res)
	return httpRes, err
}
