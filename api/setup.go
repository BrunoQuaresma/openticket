package api

import (
	"net/http"
	"reflect"
	"strings"

	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type postSetupRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (api *API) postSetup(c *gin.Context) {
	var body postSetupRequest
	c.BindJSON(body)

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	err := validate.Struct(body)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		apiErrors := make([]ValidationError, 0, len((validationErrors)))
		for _, validationError := range validationErrors {
			apiErrors = append(apiErrors, ValidationError{
				Field:     validationError.Field(),
				Validator: validationError.Tag(),
			})
		}
		c.JSON(http.StatusBadRequest, gin.H{"data": apiErrors})
		return
	}

	_, err = api.Queries.CreateUser(api.Ctx, database.CreateUserParams{
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
