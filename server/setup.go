package main

import (
	"context"
	"net/http"

	database "github.com/BrunoQuaresma/openticket/database/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postSetupBody struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func postSetup(c *gin.Context) {
	var body postSetupBody
	c.BindJSON(body)

	ctx := context.Background()

	d, err := pgxpool.New(ctx, "user=user_name dbname=db_name sslmode=disable host=localhost port=5678")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	q := database.New(d)
	_, err = q.CreateUser(ctx, database.CreateUserParams{
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
