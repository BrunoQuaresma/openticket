package api

import (
	"net/http"

	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/gin-gonic/gin"
)

type PostSetupRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Username string `json:"username" validate:"required,min=3,max=15"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (api *API) postSetup(c *gin.Context) {
	var body PostSetupRequest
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
