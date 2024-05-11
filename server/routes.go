package server

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, api *API) {
	r.GET("/health", api.getHealth)
	r.POST("/setup", api.postSetup)
}
