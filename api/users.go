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
		c.AbortWithStatusJSON(http.StatusForbidden, Response[any]{Message: "only admins can create users"})
		return
	}

	var req CreateUserRequest
	server.JSONRequest(c, &req)

	ctx := context.Background()
	tx, qtx, err := server.DBTX(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	_, err = qtx.GetUserByEmail(ctx, req.Email)
	if err == nil {
		var res Response[any]
		res.Errors = append(res.Errors, ValidationError{Field: "email", Validator: "unique"})
		res.Message = "email already in use"
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
		res.Message = "username already in use"
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
		c.AbortWithStatusJSON(http.StatusForbidden, Response[any]{Message: "only admins can delete users"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "user not found"})
		return
	}

	if user.ID == int32(id) {
		c.AbortWithStatusJSON(http.StatusForbidden, Response[any]{Message: "you can't delete yourself"})
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

type PatchUserRequest struct {
	Name     string `json:"name,omitempty" validate:"omitempty,min=3,max=50"`
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=15"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Role     string `json:"role,omitempty" validate:"omitempty,oneof=admin member"`
}

type PatchUserResponse = Response[User]

func (server *Server) patchUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "user not found"})
		return
	}

	user := server.AuthUser(c)
	if user.Role != "admin" && user.ID != int32(id) {
		c.AbortWithStatusJSON(http.StatusForbidden, Response[any]{Message: "only admins can update other users"})
		return
	}

	var req PatchUserRequest
	server.JSONRequest(c, &req)

	ctx := context.Background()
	tx, qtx, err := server.DBTX(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	u, err := qtx.GetUserByID(ctx, int32(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "user not found"})
		return
	}

	params := database.UpdateUserByIDParams{
		ID:       u.ID,
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Role:     u.Role,
	}

	if req.Name != "" {
		params.Name = req.Name
	}

	if req.Email != "" && params.Email != req.Email {
		_, err = qtx.GetUserByEmail(ctx, req.Email)
		if err == nil {
			var res Response[any]
			res.Errors = append(res.Errors, ValidationError{Field: "email", Validator: "unique"})
			res.Message = "email already in use"
			c.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}
		params.Email = req.Email
	}

	if req.Username != "" && params.Username != req.Username {
		_, err = qtx.GetUserByUsername(ctx, req.Username)
		if err == nil {
			var res Response[any]
			res.Errors = append(res.Errors, ValidationError{Field: "username", Validator: "unique"})
			res.Message = "username already in use"
			c.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}
		params.Username = req.Username
	}

	if req.Role != "" && string(u.Role) != req.Role {
		if user.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, Response[any]{Message: "only admins can update roles"})
			return
		}
		if req.Role == "member" && u.Role == database.RoleAdmin {
			countAdmins, err := qtx.CountAdmins(ctx)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			if countAdmins == 1 {
				c.AbortWithStatusJSON(http.StatusForbidden, Response[any]{Message: "can't remove the last admin"})
				return
			}
		}

		params.Role = database.Role(req.Role)
	}

	u, err = qtx.UpdateUserByID(ctx, params)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = tx.Commit(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	res := PatchUserResponse{
		Data: User{
			ID:       u.ID,
			Name:     u.Name,
			Username: u.Username,
			Email:    u.Email,
			Role:     string(u.Role),
		},
	}
	c.JSON(http.StatusOK, res)
}
