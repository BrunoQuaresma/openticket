package api

import (
	"net/http"

	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/gin-gonic/gin"
)

type postSetupRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (api *API) postSetup(c *gin.Context) {
	var body postSetupRequest
	api.BodyAsJSON(&body, c)

	_, err := api.Queries.CreateUser(api.Context, database.CreateUserParams{
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
