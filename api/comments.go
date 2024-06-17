package api

import (
	"net/http"
	"strconv"

	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=10"`
	ReplyTo int32  `json:"reply_to,omitempty" validate:"number,omitempty"`
}

type Comment struct {
	ID        int32  `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	ReplyTo   int32  `json:"reply_to,omitempty"`
	CreatedBy User   `json:"created_by"`
}

type CreateCommentResponse = Response[Comment]

func (server *Server) createComment(c *gin.Context) {
	user := server.AuthUser(c)

	ticketId, err := strconv.ParseUint(c.Param("ticketId"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "ticket not found"})
		return
	}

	var req CreateCommentRequest
	server.jsonReq(c, &req)

	ctx := c.Request.Context()
	newComment, err := server.db.queries.CreateComment(ctx, database.CreateCommentParams{
		Content:  req.Content,
		TicketID: int32(ticketId),
		UserID:   user.ID,
		ReplyTo:  pgtype.Int4{Int32: req.ReplyTo, Valid: req.ReplyTo != 0},
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response[any]{Message: "failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, CreateCommentResponse{
		Data: Comment{
			ID:        newComment.ID,
			Content:   newComment.Content,
			CreatedAt: newComment.CreatedAt.Time.UTC().String(),
			ReplyTo:   newComment.ReplyTo.Int32,
			CreatedBy: User{
				ID:       user.ID,
				Username: user.Username,
				Name:     user.Name,
				Email:    user.Email,
				Role:     string(user.Role),
			},
		},
	})
}
