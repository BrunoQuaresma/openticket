package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	req := gin.Default()
	req.POST("/setup", postSetup)
	req.Run()
}
