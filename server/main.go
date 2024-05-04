package main

import (
	"context"
	"os"

	database "github.com/BrunoQuaresma/openticket/database/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type API struct {
	Ctx     *context.Context
	Queries *database.Queries
}

func main() {
	ctx := context.Background()

	d, err := pgxpool.New(ctx, os.Getenv("POSTGRES_DB_URL"))
	if err != nil {
		panic("Error connecting to the database. Error: " + err.Error())
	}

	api := &API{
		Ctx:     &ctx,
		Queries: database.New(d),
	}

	req := gin.Default()
	req.POST("/setup", api.postSetup)
	req.Run()
}
