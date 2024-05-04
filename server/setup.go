package main

import (
	"net/http"

	database "github.com/BrunoQuaresma/openticket/database/models"
	"github.com/gin-gonic/gin"
)

type postSetupBody struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (api *API) postSetup(c *gin.Context) {
	var body postSetupBody
	c.BindJSON(body)

	_, err := api.Queries.CreateUser(api.Ctx, database.CreateUserParams{
		Name:     body.Name,
		Username: body.Username,
		Email:    body.Email,
		Hash:     body.Password,
	})

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
