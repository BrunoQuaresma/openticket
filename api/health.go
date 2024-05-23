package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) health(c *gin.Context) {
	c.Status(http.StatusOK)
}
