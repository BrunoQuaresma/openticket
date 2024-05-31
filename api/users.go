package api

import (
	"context"
	"net/http"
	"strconv"

	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Username string `json:"username" validate:"required,min=3,max=15"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=admin member"`
}

type User struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type CreateUserResponse = Response[User]

func (server *Server) createUser(c *gin.Context) {
	user := server.AuthUser(c)
	if user.Role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var req CreateUserRequest
	server.JSONRequest(c, &req)

	ctx := context.Background()
	dbQueries := server.DBQueries()
	tx, err := server.BeginTX(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback(ctx)
	qtx := dbQueries.WithTx(tx)

	_, err = qtx.GetUserByEmail(ctx, req.Email)
	if err == nil {
		var res Response[any]
		res.Errors = append(res.Errors, ValidationError{Field: "email", Validator: "unique"})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	if err != pgx.ErrNoRows {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	_, err = qtx.GetUserByUsername(ctx, req.Username)
	if err == nil {
		var res Response[any]
		res.Errors = append(res.Errors, ValidationError{Field: "username", Validator: "unique"})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	if err != pgx.ErrNoRows {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	h, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	u, err := qtx.CreateUser(ctx, database.CreateUserParams{
		Name:         req.Name,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(h),
		Role:         database.Role(req.Role),
	})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	tx.Commit(ctx)

	res := CreateUserResponse{
		Data: User{
			ID:       u.ID,
			Name:     u.Name,
			Username: u.Username,
			Email:    u.Email,
			Role:     string(u.Role),
		},
	}

	c.JSON(http.StatusCreated, res)
}

func (server *Server) deleteUser(c *gin.Context) {
	user := server.AuthUser(c)
	if user.Role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if user.ID == int32(id) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	ctx := context.Background()
	dbQueries := server.DBQueries()
	err = dbQueries.DeleteUserByID(ctx, int32(id))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
