package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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
	server.jsonReq(c, &req)

	var (
		ticket database.Ticket
		err    error
	)
	err = server.db.tx(func(ctx context.Context, qtx *database.Queries, _ pgx.Tx) error {
		ticket, err = qtx.CreateTicket(ctx, database.CreateTicketParams{
			Title:       req.Title,
			Description: req.Description,
			CreatedBy:   user.ID,
		})
		if err != nil {
			return err
		}

		if len(req.Labels) > 0 {
			for _, labelName := range req.Labels {
				label, err := qtx.GetLabelByName(ctx, labelName)
				if err == pgx.ErrNoRows {
					label, err = qtx.CreateLabelIfNotExists(ctx, database.CreateLabelIfNotExistsParams{
						Name:      labelName,
						CreatedBy: user.ID,
					})
					if err != nil {
						return err
					}
				} else if err != nil {
					return err
				}

				err = qtx.AssignLabelToTicketIfNotAssigned(ctx, database.AssignLabelToTicketIfNotAssignedParams{
					TicketID: ticket.ID,
					Name:     label.Name,
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

	var ticketRows []database.GetTicketsByIDsRow
	err := server.db.tx(func(ctx context.Context, qtx *database.Queries, tx pgx.Tx) error {
		baseSelect := "SELECT tickets.id FROM tickets " +
			"JOIN ticket_labels ON tickets.id = ticket_labels.ticket_id " +
			"JOIN labels ON ticket_labels.label_id = labels.id " +
			"JOIN users ON tickets.created_by = users.id "
		selects := []string{}
		if len(tags) > 0 {
			for _, tag := range tags {
				filterQuery := baseSelect + "WHERE "
				switch tag.Key {
				case "title":
					filterQuery += "tickets.title ILIKE '%" + strings.Join(tag.Values, "%' AND tickets.title ILIKE '%") + "%' "
				case "label":
					filterQuery += "labels.name IN ('" + strings.Join(tag.Values, "', '") + "') "
				default:
					return errors.New("invalid tag key: " + tag.Key)
				}
				selects = append(selects, filterQuery)
			}
		}
		selectQuery := strings.Join(selects, "INTERSECT ") + ";"

		rows, err := tx.Query(ctx, selectQuery)
		if err != nil {
			return err
		}
		defer rows.Close()
		results, err := pgx.CollectRows(rows, pgx.RowToStructByName[struct{ ID int32 }])
		if err != nil {
			return err
		}
		ids := make([]int32, len(results))
		for i, result := range results {
			ids[i] = result.ID
		}
		ticketRows, err = qtx.GetTicketsByIDs(ctx, ids)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tickets := make([]Ticket, len(ticketRows))
	for i, ticket := range ticketRows {
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

type TicketNotFoundError struct{}

func (e TicketNotFoundError) Error() string {
	return "ticket not found"
}

func (server *Server) deleteTicket(c *gin.Context) {
	user := server.AuthUser(c)

	ticketId, err := strconv.ParseInt(c.Param("ticketId"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "ticket not found"})
		return
	}

	err = server.db.tx(func(ctx context.Context, qtx *database.Queries, _ pgx.Tx) error {
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
	user := server.AuthUser(c)

	ticketId, err := strconv.ParseInt(c.Param("ticketId"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "ticket not found"})
		return
	}

	var req PatchTicketRequest
	server.jsonReq(c, &req)

	var (
		updatedTicket database.Ticket
		createdBy     database.User
	)
	err = server.db.tx(func(ctx context.Context, qtx *database.Queries, _ pgx.Tx) error {
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
		if req.Description != "" {
			ticket.Description = req.Description
		}

		if len(req.Labels) > 0 {
			for _, labelName := range req.Labels {
				label, err := qtx.GetLabelByName(ctx, labelName)
				if err == pgx.ErrNoRows {
					label, err = qtx.CreateLabelIfNotExists(ctx, database.CreateLabelIfNotExistsParams{
						Name:      labelName,
						CreatedBy: user.ID,
					})
					if err != nil {
						return err
					}
				} else if err != nil {
					return err
				}

				err = qtx.AssignLabelToTicketIfNotAssigned(ctx, database.AssignLabelToTicketIfNotAssignedParams{
					TicketID: ticket.ID,
					Name:     label.Name,
				})
				if err != nil {
					return err
				}
			}
		}

		updatedTicket, err = qtx.UpdateTicketByID(ctx, database.UpdateTicketByIDParams{
			ID:          ticket.ID,
			Title:       ticket.Title,
			Description: ticket.Description,
		})
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
				ID:          updatedTicket.ID,
				Title:       updatedTicket.Title,
				Description: updatedTicket.Description,
				Labels:      req.Labels,
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
