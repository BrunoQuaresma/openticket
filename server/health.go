package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *API) getHealth(c *gin.Context) {
	c.Status(http.StatusOK)
}
