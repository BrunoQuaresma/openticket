package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Health struct {
	Setup    bool `json:"setup"`
	Database bool `json:"database"`
}

type HealthResponse = Response[Health]

func (server *Server) health(c *gin.Context) {
	hasFirstUser, err := server.db.queries.HasFirstUser(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, HealthResponse{
		Data: Health{
			Setup:    hasFirstUser,
			Database: true,
		},
	})
}
