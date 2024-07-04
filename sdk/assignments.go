package sdk

import (
	"fmt"
	"net/http"

	"github.com/BrunoQuaresma/openticket/api"
)

func (c *Client) CreateAssignment(ticketId int32, userId int32, res *api.CreateAssignmentResponse) (*http.Response, error) {
	httpRes, err := c.post(
		"/tickets/"+fmt.Sprint(ticketId)+"/assignments",
		api.CreateAssignmentRequest{UserID: userId},
		res,
	)
	return httpRes, err
}

func (c *Client) DeleteAssignment(assignmentId int32) (*http.Response, error) {
	return c.delete("/tickets/1/assignments/" + fmt.Sprint(assignmentId))
}
