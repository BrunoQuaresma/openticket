package api

import (
	"context"
	"net/http"
	"strconv"

	sqlc "github.com/BrunoQuaresma/openticket/api/database/sqlc"
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
	authUser := server.AuthUserFromContext(c)
	if authUser.Role != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, Response[any]{Message: "only admins can create users"})
		return
	}

	var req CreateUserRequest
	server.jsonReq(c, &req)

	var user sqlc.User
	err := server.db.TX(func(ctx context.Context, qtx *sqlc.Queries, _ pgx.Tx) error {
		_, err := qtx.GetUserByEmail(ctx, req.Email)
		if err == nil {
			return EmailAlreadyInUseError{}
		}
		if err != pgx.ErrNoRows {
			return err
		}

		_, err = qtx.GetUserByUsername(ctx, req.Username)
		if err == nil {
			return UsernameAlreadyInUseError{}
		}
		if err != pgx.ErrNoRows {
			return err
		}

		h, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user, err = qtx.CreateUser(ctx, sqlc.CreateUserParams{
			Name:         req.Name,
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: string(h),
			Role:         sqlc.Role(req.Role),
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		switch err.(type) {
		case EmailAlreadyInUseError:
			c.AbortWithStatusJSON(http.StatusBadRequest, Response[any]{
				Message: err.Error(),
				Errors: []ValidationError{
					{Field: "email", Validator: "unique"},
				},
			})
		case UsernameAlreadyInUseError:
			c.AbortWithStatusJSON(http.StatusBadRequest, Response[any]{
				Message: err.Error(),
				Errors: []ValidationError{
					{Field: "username", Validator: "unique"},
				},
			})
		default:
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	res := CreateUserResponse{
		Data: User{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			Role:     string(user.Role),
		},
	}

	c.JSON(http.StatusCreated, res)
}

func (server *Server) deleteUser(c *gin.Context) {
	user := server.AuthUserFromContext(c)
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
	err = server.db.Queries().DeleteUserByID(ctx, int32(id))
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

	authUser := server.AuthUserFromContext(c)
	if authUser.Role != "admin" && authUser.ID != int32(id) {
		c.AbortWithStatusJSON(http.StatusForbidden, Response[any]{
			Message: "only admins can update other users",
		})
		return
	}

	var req PatchUserRequest
	server.jsonReq(c, &req)

	var updatedUser sqlc.User
	err = server.db.TX(func(ctx context.Context, qtx *sqlc.Queries, _ pgx.Tx) error {
		u, err := qtx.GetUserByID(ctx, int32(id))
		if err != nil {
			return UserNotFoundError{}
		}

		params := sqlc.UpdateUserByIDParams{
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
				return EmailAlreadyInUseError{}
			}
			params.Email = req.Email
		}

		if req.Username != "" && params.Username != req.Username {
			_, err = qtx.GetUserByUsername(ctx, req.Username)
			if err == nil {
				return UsernameAlreadyInUseError{}
			}
			params.Username = req.Username
		}

		if req.Role != "" && string(u.Role) != req.Role {
			if authUser.Role != "admin" {
				return PermissionDeniedError{Message: "only admins can update roles"}
			}
			if req.Role == "member" && u.Role == sqlc.RoleAdmin {
				countAdmins, err := qtx.CountAdmins(ctx)
				if err != nil {
					return err
				}
				if countAdmins == 1 {
					return PermissionDeniedError{Message: "can't remove the last admin"}
				}
			}

			params.Role = sqlc.Role(req.Role)
		}

		updatedUser, err = qtx.UpdateUserByID(ctx, params)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		switch err.(type) {
		case UserNotFoundError:
			c.AbortWithStatusJSON(http.StatusNotFound, Response[any]{Message: "user not found"})
		case EmailAlreadyInUseError:
			c.AbortWithStatusJSON(http.StatusBadRequest, Response[any]{
				Message: err.Error(),
				Errors: []ValidationError{
					{Field: "email", Validator: "unique"},
				},
			})
		case UsernameAlreadyInUseError:
			c.AbortWithStatusJSON(http.StatusBadRequest, Response[any]{
				Message: err.Error(),
				Errors: []ValidationError{
					{Field: "username", Validator: "unique"},
				},
			})
		case PermissionDeniedError:
			c.AbortWithStatusJSON(http.StatusForbidden, Response[any]{
				Message: err.Error(),
			})
		default:
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	res := PatchUserResponse{
		Data: User{
			ID:       updatedUser.ID,
			Name:     updatedUser.Name,
			Username: updatedUser.Username,
			Email:    updatedUser.Email,
			Role:     string(updatedUser.Role),
		},
	}
	c.JSON(http.StatusOK, res)
}

type EmailAlreadyInUseError struct{}

func (e EmailAlreadyInUseError) Error() string {
	return "email already in use"
}

type UsernameAlreadyInUseError struct{}

func (e UsernameAlreadyInUseError) Error() string {
	return "username already in use"
}

type UserNotFoundError struct{}

func (e UserNotFoundError) Error() string {
	return "user not found"
}
