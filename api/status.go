package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Status struct {
	Setup bool  `json:"setup"`
	User  *User `json:"user,omitempty"`
}

type StatusResponse = Response[Status]

func (s *Server) status(c *gin.Context) {
	authUser, err := s.AuthUser(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	hasFirstUser, err := s.db.Queries().HasFirstUser(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var user *User
	if authUser != nil {
		user = &User{
			ID:       authUser.ID,
			Username: authUser.Username,
			Email:    authUser.Email,
			Name:     authUser.Name,
		}
	}

	c.JSON(200, StatusResponse{
		Data: Status{
			Setup: hasFirstUser,
			User:  user,
		},
	})
}
