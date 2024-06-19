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

func (c *Client) DeleteComment(ticketId int32, commentId int32) (*http.Response, error) {
	httpRes, err := c.Delete("/tickets/" + fmt.Sprint(ticketId) + "/comments/" + fmt.Sprint(commentId))
	return httpRes, err
}
