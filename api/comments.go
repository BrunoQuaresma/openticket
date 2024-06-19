package api

import (
	"context"
	"net/http"
	"strconv"

	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

type CommentNotFoundError struct{}

func (e CommentNotFoundError) Error() string {
	return "comment not found"
}

func (server *Server) deleteComment(c *gin.Context) {
	user := server.AuthUser(c)

	commentId, err := strconv.ParseInt(c.Param("commentId"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "comment not found"})
		return
	}

	err = server.db.tx(func(ctx context.Context, qtx *database.Queries, _ pgx.Tx) error {
		comment, err := qtx.GetCommentByID(ctx, int32(commentId))
		if err != nil {
			return CommentNotFoundError{}
		}

		if comment.UserID == user.ID || user.Role == "admin" {
			return qtx.DeleteComment(ctx, int32(commentId))
		}

		return PermissionDeniedError{Message: "only admins and the comment's author can delete comments"}
	})

	switch err.(type) {
	case nil:
		c.Status(http.StatusNoContent)
	case CommentNotFoundError:
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "comment not found"})
	case PermissionDeniedError:
		c.AbortWithStatusJSON(http.StatusForbidden, Response[any]{Message: err.Error()})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response[any]{Message: "failed to delete comment"})
	}
}

type PatchCommentRequest struct {
	Content string `json:"content" validate:"required,min=10"`
}

type PatchCommentResponse = Response[Comment]

func (server *Server) patchComment(c *gin.Context) {
	user := server.AuthUser(c)

	commentId, err := strconv.ParseInt(c.Param("commentId"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "comment not found"})
		return
	}

	var req PatchCommentRequest
	server.jsonReq(c, &req)

	var (
		updatedComment database.Comment
		commentOwner   database.User
	)
	err = server.db.tx(func(ctx context.Context, qtx *database.Queries, _ pgx.Tx) error {
		comment, err := qtx.GetCommentByID(ctx, int32(commentId))
		if err != nil {
			return CommentNotFoundError{}
		}

		if comment.UserID == user.ID || user.Role == "admin" {
			updatedComment, err = qtx.UpdateCommentByID(ctx, database.UpdateCommentByIDParams{
				ID:      comment.ID,
				Content: req.Content,
			})
			if err != nil {
				return err
			}
			commentOwner, err = qtx.GetUserByID(ctx, comment.UserID)
			return err
		}

		return PermissionDeniedError{Message: "only the comment's author or admins can edit comments"}
	})

	switch err.(type) {
	case nil:
		c.JSON(http.StatusOK, PatchCommentResponse{
			Data: Comment{
				ID:        updatedComment.ID,
				Content:   updatedComment.Content,
				CreatedAt: updatedComment.CreatedAt.Time.UTC().String(),
				ReplyTo:   updatedComment.ReplyTo.Int32,
				CreatedBy: User{
					ID:       commentOwner.ID,
					Username: commentOwner.Username,
					Name:     commentOwner.Name,
					Email:    commentOwner.Email,
					Role:     string(commentOwner.Role),
				},
			},
		})
	case CommentNotFoundError:
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "comment not found"})
	case PermissionDeniedError:
		c.AbortWithStatusJSON(http.StatusForbidden, Response[any]{Message: err.Error()})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response[any]{Message: "failed to update comment"})
	}
}
