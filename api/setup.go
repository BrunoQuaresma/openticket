package api

import (
	"context"
	"net/http"

	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/gin-gonic/gin"
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

func (server *Server) setup(c *gin.Context) {
	var req SetupRequest
	server.jsonReq(c, &req)

	ctx := context.Background()
	tx, qtx, err := server.DBTX(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	count, err := qtx.CountUsers(ctx)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if count > 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	h, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user, err := qtx.CreateUser(ctx, database.CreateUserParams{
		Name:         req.Name,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(h),
		Role:         "admin",
	})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	tx.Commit(ctx)

	c.JSON(http.StatusOK, SetupResponse{
		Data: Setup{
			ID:   user.ID,
			Role: string(user.Role),
		},
	})
}
