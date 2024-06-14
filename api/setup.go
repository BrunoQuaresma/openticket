package api

import (
	"context"
	"net/http"

	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type SetupRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Username string `json:"username" validate:"required,min=3,max=15"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type Setup struct {
	ID   int32  `json:"id"`
	Role string `json:"role"`
}

type SetupResponse = Response[Setup]

type SetupAlreadyDoneError struct{}

func (e SetupAlreadyDoneError) Error() string {
	return "setup is already done"
}

func (server *Server) setup(c *gin.Context) {
	var req SetupRequest
	server.jsonReq(c, &req)

	var user database.User
	err := server.db.tx(func(ctx context.Context, qtx *database.Queries, _ pgx.Tx) error {
		count, err := qtx.CountUsers(ctx)
		if err != nil {
			return err
		}
		if count > 0 {
			return SetupAlreadyDoneError{}
		}

		h, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user, err = qtx.CreateUser(ctx, database.CreateUserParams{
			Name:         req.Name,
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: string(h),
			Role:         "admin",
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		switch err.(type) {
		case SetupAlreadyDoneError:
			c.AbortWithStatusJSON(http.StatusForbidden, Response[any]{Message: err.Error()})
		default:
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, SetupResponse{
		Data: Setup{
			ID:   user.ID,
			Role: string(user.Role),
		},
	})
}
