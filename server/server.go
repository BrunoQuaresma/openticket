package server

import (
	"context"
	"net/http"

	database "github.com/BrunoQuaresma/openticket/database/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type API struct {
	Ctx     context.Context
	Queries *database.Queries
}

type Options struct {
	DatabaseURL string
	Debug       bool
}

func Start(options Options) *http.Server {
	ctx := context.Background()

	d, err := pgxpool.New(ctx, options.DatabaseURL)
	if err != nil {
		panic("error connecting to the database. " + err.Error())
	}

	api := API{
		Ctx:     ctx,
		Queries: database.New(d),
	}
	if !options.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	SetupRoutes(r, &api)

	server := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}
	go func() {
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic("error starting server. " + err.Error())
		}
	}()
	ready := make(chan bool, 1)
	go func() {
		var res *http.Response

		for res == nil || res.StatusCode != http.StatusOK {
			res, err = http.Get("http://localhost:3000/health")
			if err != nil {
				panic("error getting health check. " + err.Error())
			}
		}

		ready <- true
	}()
	<-ready

	return server
}
