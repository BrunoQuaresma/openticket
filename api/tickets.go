package api

import (
	"context"
	"net/http"

	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/gin-gonic/gin"
)

type CreateTicketRequest struct {
	Title       string   `json:"title" validate:"required,min=3,max=70"`
	Description string   `json:"description" validate:"required,min=10"`
	Labels      []string `json:"labels,omitempty"`
}

type Ticket struct {
	ID          int32    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Labels      []string `json:"labels"`
	CreatedBy   User     `json:"created_by"`
}

type CreateTicketResponse = Response[Ticket]

func (server *Server) createTicket(c *gin.Context) {
	user := server.AuthUser(c)

	var req CreateTicketRequest
	server.JSONRequest(c, &req)

	ctx := context.Background()
	tx, qtx, err := server.DBTX(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ticket, err := qtx.CreateTicket(ctx, database.CreateTicketParams{
		Title:       req.Title,
		Description: req.Description,
		CreatedBy:   user.ID,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(req.Labels) > 0 {
		for _, labelName := range req.Labels {
			label, err := qtx.CreateLabel(ctx, database.CreateLabelParams{
				Name:      labelName,
				CreatedBy: user.ID,
			})
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			err = qtx.AssignLabelToTicket(ctx, database.AssignLabelToTicketParams{
				TicketID: ticket.ID,
				LabelID:  label.ID,
			})
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, CreateTicketResponse{
		Data: Ticket{
			ID:          ticket.ID,
			Title:       ticket.Title,
			Description: ticket.Description,
			Labels:      req.Labels,
			CreatedBy: User{
				ID:       user.ID,
				Name:     user.Name,
				Username: user.Username,
				Email:    user.Email,
				Role:     string(user.Role),
			},
		},
	})
}

type TicketsResponse = Response[[]Ticket]

func (server *Server) tickets(c *gin.Context) {
	ctx := context.Background()
	dbQueries := server.DBQueries()
	result, err := dbQueries.GetTickets(ctx, []string{})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tickets := make([]Ticket, len(result))
	for i, ticket := range result {
		tickets[i] = Ticket{
			ID:          ticket.ID,
			Title:       ticket.Title,
			Description: ticket.Description,
			Labels:      ticket.Labels,
			CreatedBy: User{
				ID:       ticket.User.ID,
				Name:     ticket.User.Name,
				Username: ticket.User.Username,
				Email:    ticket.User.Email,
				Role:     string(ticket.User.Role),
			},
		}
	}
	c.JSON(http.StatusOK, TicketsResponse{
		Data: tickets,
	})
}
