package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) createTicket(c *gin.Context) {
	c.Status(http.StatusOK)
}
