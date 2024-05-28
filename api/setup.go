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

func (server *Server) setup(c *gin.Context) {
	var req SetupRequest
	server.ParseJSONRequest(c, &req)

	ctx := context.Background()
	dbQueries := server.DBQueries()
	tx, err := server.BeginTX(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback(ctx)

	qtx := dbQueries.WithTx(tx)
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

	_, err = qtx.CreateUser(ctx, database.CreateUserParams{
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

	c.Status(http.StatusOK)
}
