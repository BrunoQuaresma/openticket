package api

import (
	"net/http"
	"strconv"

	sqlc "github.com/BrunoQuaresma/openticket/api/database/sqlc"
	"github.com/gin-gonic/gin"
)

type CreateAssignmentRequest struct {
	UserID int32 `json:"user_id" validate:"required,number"`
}

type Assignment struct {
	ID       int32 `json:"id"`
	TicketID int32 `json:"ticket_id"`
	UserID   int32 `json:"user_id"`
}

type CreateAssignmentResponse = Response[Assignment]

func (server *Server) createAssignment(c *gin.Context) {
	user := server.AuthUser(c)

	var req CreateAssignmentRequest
	server.jsonReq(c, &req)

	ticketId, err := strconv.ParseUint(c.Param("ticketId"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "ticket not found"})
		return
	}

	assignment, err := server.db.Queries().CreateAssignment(c.Request.Context(), sqlc.CreateAssignmentParams{
		TicketID:   int32(ticketId),
		UserID:     req.UserID,
		AssignedBy: user.ID,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response[any]{Message: "failed to create assignment"})
		return
	}

	c.JSON(http.StatusCreated, CreateAssignmentResponse{
		Data: Assignment{
			TicketID: assignment.TicketID,
			UserID:   assignment.UserID,
			ID:       assignment.ID,
		},
	})
}

func (server *Server) deleteAssignment(c *gin.Context) {
	assignmentId, err := strconv.ParseUint(c.Param("assignmentId"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "assignment not found"})
		return
	}

	err = server.db.Queries().DeleteAssignment(c.Request.Context(), int32(assignmentId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response[any]{Message: "failed to delete assignment"})
		return
	}

	c.JSON(http.StatusOK, Response[any]{Message: "assignment deleted"})
}
