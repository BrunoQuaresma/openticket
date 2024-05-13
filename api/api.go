package api

import (
	"context"
	"fmt"
	"net/http"

	database "github.com/BrunoQuaresma/openticket/api/database/gen"
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
	Port        int
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
	setupRoutes(r, &api)

	server := &http.Server{
		Addr:    ":" + fmt.Sprint(options.Port),
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
			res, err = http.Get("http://localhost" + server.Addr + "/health")
			if err != nil {
				panic("error getting health check. " + err.Error())
			}
		}

		ready <- true
	}()
	<-ready

	return server
}

func setupRoutes(r *gin.Engine, api *API) {
	r.GET("/health", api.getHealth)
	r.POST("/setup", api.postSetup)
}
