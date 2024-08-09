package api

import (
	"context"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	sqlc "github.com/BrunoQuaresma/openticket/api/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type CreateTicketRequest struct {
	Title       string   `json:"title" validate:"required,min=3,max=70"`
	Description string   `json:"description" validate:"required,min=10"`
	Labels      []string `json:"labels,omitempty"`
}

type Ticket struct {
	ID        int32    `json:"id"`
	Title     string   `json:"title"`
	Status    string   `json:"status"`
	Labels    []string `json:"labels"`
	CreatedBy User     `json:"created_by"`
	CreatedAt string   `json:"created_at"`
}

type CreateTicketResponse = Response[Ticket]

func (server *Server) createTicket(c *gin.Context) {
	user := server.AuthUserFromContext(c)

	var req CreateTicketRequest
	server.jsonReq(c, &req)

	var (
		ticket sqlc.Ticket
		err    error
	)
	err = server.db.TX(func(ctx context.Context, qtx *sqlc.Queries, _ pgx.Tx) error {
		ticket, err = qtx.CreateTicket(ctx, sqlc.CreateTicketParams{
			Title:     req.Title,
			CreatedBy: user.ID,
		})
		if err != nil {
			return err
		}

		_, err = qtx.CreateComment(ctx, sqlc.CreateCommentParams{
			Content:  req.Description,
			TicketID: ticket.ID,
			UserID:   user.ID,
		})
		if err != nil {
			return err
		}

		if len(req.Labels) > 0 {
			for _, labelName := range req.Labels {
				err = qtx.AssignLabelToTicket(ctx, sqlc.AssignLabelToTicketParams{
					TicketID:  ticket.ID,
					LabelName: labelName,
				})
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, CreateTicketResponse{
		Data: Ticket{
			ID:        ticket.ID,
			Title:     ticket.Title,
			Status:    string(ticket.Status),
			Labels:    req.Labels,
			CreatedAt: ticket.CreatedAt.Time.Format(time.RFC3339),
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

type Tag struct {
	Key    string
	Values []string
}

func (server *Server) tickets(c *gin.Context) {
	var tags []Tag
	q, hasQuery := c.GetQuery("q")
	if hasQuery {
		sentences := strings.Split(q, " ")
		for _, sentence := range sentences {
			words := strings.Split(sentence, ":")
			if len(words) == 1 {
				words = []string{"title", words[0]}
			} else if len(words) != 2 {
				c.AbortWithStatusJSON(http.StatusBadRequest, Response[any]{
					Message: "Invalid query",
					Errors: []ValidationError{
						{Field: "q", Validator: "search"},
					},
				})
				return
			}
			key := words[0]
			values := strings.Split(words[1], ",")
			tags = append(tags, Tag{
				Key:    key,
				Values: values,
			})
		}
	}

	var (
		title  string
		labels []string
	)

	if len(tags) > 0 {
		for _, tag := range tags {
			switch tag.Key {
			case "title":
				title = tag.Values[0]
			case "label":
				labels = tag.Values
			}
		}
	}

	ticketRows, err := server.db.Queries().GetTickets(c, sqlc.GetTicketsParams{
		Labels: labels,
		Title:  title,
	})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tickets := make([]Ticket, len(ticketRows))
	for i, ticket := range ticketRows {
		tickets[i] = Ticket{
			ID:        ticket.ID,
			Title:     ticket.Title,
			Status:    string(ticket.Status),
			Labels:    ticket.Labels,
			CreatedAt: ticket.CreatedAt.Time.Format(time.RFC3339),
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

type TicketNotFoundError struct{}

func (e TicketNotFoundError) Error() string {
	return "ticket not found"
}

func (server *Server) deleteTicket(c *gin.Context) {
	user := server.AuthUserFromContext(c)

	ticketId, err := strconv.ParseInt(c.Param("ticketId"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "ticket not found"})
		return
	}

	err = server.db.TX(func(ctx context.Context, qtx *sqlc.Queries, _ pgx.Tx) error {
		ticket, err := qtx.GetTicketByID(ctx, int32(ticketId))
		if err != nil {
			return TicketNotFoundError{}
		}

		if ticket.CreatedBy == user.ID || user.Role == "admin" {
			return qtx.DeleteTicketByID(ctx, int32(ticketId))
		}

		return PermissionDeniedError{Message: "only admins and the ticket's creator can delete tickets"}
	})

	switch err.(type) {
	case nil:
		c.Status(http.StatusNoContent)
	case TicketNotFoundError:
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "ticket not found"})
	case PermissionDeniedError:
		c.AbortWithStatusJSON(http.StatusForbidden, Response[any]{Message: err.Error()})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response[any]{Message: "failed to delete ticket"})
	}
}

type PatchTicketRequest struct {
	Title       string   `json:"title,omitempty" validate:"omitempty,min=3,max=70"`
	Description string   `json:"description,omitempty" validate:"omitempty,min=10"`
	Labels      []string `json:"labels,omitempty"`
}

type PatchTicketResponse = Response[Ticket]

func (server *Server) patchTicket(c *gin.Context) {
	user := server.AuthUserFromContext(c)

	ticketId, err := strconv.ParseInt(c.Param("ticketId"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "ticket not found"})
		return
	}

	var req PatchTicketRequest
	server.jsonReq(c, &req)

	var (
		updatedTicket sqlc.GetTicketByIDRow
		createdBy     sqlc.User
	)
	err = server.db.TX(func(ctx context.Context, qtx *sqlc.Queries, _ pgx.Tx) error {
		ticket, err := qtx.GetTicketByID(ctx, int32(ticketId))
		if err != nil {
			return TicketNotFoundError{}
		}

		if ticket.CreatedBy != user.ID && user.Role != "admin" {
			return PermissionDeniedError{Message: "only admins and the ticket's creator can update tickets"}
		}

		if req.Title != "" {
			ticket.Title = req.Title
		}

		var labelsToUnassign []string
		for _, oldLabel := range ticket.Labels {
			if !slices.Contains(req.Labels, oldLabel) {
				labelsToUnassign = append(labelsToUnassign, oldLabel)
			}
		}

		var labelsToAssign []string
		for _, newLabel := range req.Labels {
			if !slices.Contains(ticket.Labels, newLabel) {
				labelsToAssign = append(labelsToAssign, newLabel)
			}
		}

		for _, labelName := range labelsToUnassign {
			err := qtx.UnassignLabelFromTicket(ctx, sqlc.UnassignLabelFromTicketParams{
				TicketID:  ticket.ID,
				LabelName: labelName,
			})
			if err != nil {
				return err
			}
		}

		for _, labelName := range labelsToAssign {
			err := qtx.AssignLabelToTicket(ctx, sqlc.AssignLabelToTicketParams{
				TicketID:  ticket.ID,
				LabelName: labelName,
			})
			if err != nil {
				return err
			}
		}

		_, err = qtx.UpdateTicketByID(ctx, sqlc.UpdateTicketByIDParams{
			ID:    ticket.ID,
			Title: ticket.Title,
		})
		if err != nil {
			return err
		}

		updatedTicket, err = qtx.GetTicketByID(ctx, ticket.ID)
		if err != nil {
			return err
		}

		createdBy, err = qtx.GetUserByID(ctx, ticket.CreatedBy)
		return err
	})

	switch err.(type) {
	case nil:
		c.JSON(http.StatusOK, PatchTicketResponse{
			Data: Ticket{
				ID:        updatedTicket.ID,
				Title:     updatedTicket.Title,
				Status:    string(updatedTicket.Status),
				Labels:    updatedTicket.Labels,
				CreatedAt: updatedTicket.CreatedAt.Time.Format(time.RFC3339),
				CreatedBy: User{
					ID:       createdBy.ID,
					Name:     createdBy.Name,
					Username: createdBy.Username,
					Email:    createdBy.Email,
					Role:     string(createdBy.Role),
				},
			},
		})
	case TicketNotFoundError:
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "ticket not found"})
	case PermissionDeniedError:
		c.AbortWithStatusJSON(http.StatusForbidden, Response[any]{Message: err.Error()})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response[any]{Message: "failed to update ticket"})
	}
}

type TicketResponse = Response[Ticket]

func (server *Server) ticket(c *gin.Context) {
	ticketId, err := strconv.ParseInt(c.Param("ticketId"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "ticket not found"})
		return
	}

	var ticketRow sqlc.GetTicketByIDRow
	err = server.db.TX(func(ctx context.Context, qtx *sqlc.Queries, _ pgx.Tx) error {
		ticket, err := qtx.GetTicketByID(ctx, int32(ticketId))
		if err != nil {
			return TicketNotFoundError{}
		}

		ticketRow = ticket
		return nil
	})

	switch err.(type) {
	case nil:
		c.JSON(http.StatusOK, PatchTicketResponse{
			Data: Ticket{
				ID:        ticketRow.ID,
				Title:     ticketRow.Title,
				Status:    string(ticketRow.Status),
				Labels:    ticketRow.Labels,
				CreatedAt: ticketRow.CreatedAt.Time.Format(time.RFC3339),
				CreatedBy: User{
					ID:       ticketRow.User.ID,
					Name:     ticketRow.User.Name,
					Username: ticketRow.User.Username,
					Email:    ticketRow.User.Email,
					Role:     string(ticketRow.User.Role),
				},
			},
		})
	case TicketNotFoundError:
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "ticket not found"})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response[any]{Message: "failed to get ticket"})
	}
}

type PatchTicketStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=open closed"`
}

type PatchTicketStatusResponse = Response[Ticket]

func (server *Server) patchTicketStatus(c *gin.Context) {
	ticketId, err := strconv.ParseInt(c.Param("ticketId"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "ticket not found"})
		return
	}

	var req PatchTicketStatusRequest
	server.jsonReq(c, &req)

	_, err = server.db.Queries().UpdateTicketStatusByID(
		c.Request.Context(),
		sqlc.UpdateTicketStatusByIDParams{
			ID:     int32(ticketId),
			Status: sqlc.TicketStatus(req.Status),
		},
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response[any]{Message: "failed to update ticket status"})
		return
	}
	updatedTicket, err := server.db.Queries().GetTicketByID(c.Request.Context(), int32(ticketId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response[any]{Message: "failed to get updated ticket"})
		return
	}

	c.JSON(http.StatusOK, PatchTicketStatusResponse{
		Data: Ticket{
			ID:        updatedTicket.ID,
			Title:     updatedTicket.Title,
			Labels:    updatedTicket.Labels,
			Status:    string(updatedTicket.Status),
			CreatedAt: updatedTicket.CreatedAt.Time.Format(time.RFC3339),
			CreatedBy: User{
				ID:       updatedTicket.User.ID,
				Name:     updatedTicket.User.Name,
				Username: updatedTicket.User.Username,
				Email:    updatedTicket.User.Email,
				Role:     string(updatedTicket.User.Role),
			},
		},
	})
}
