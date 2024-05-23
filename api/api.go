package api

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	sqlc "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type API struct {
	Context    context.Context
	Queries    *sqlc.Queries
	Database   *pgxpool.Pool
	validate   *validator.Validate
	httpServer *http.Server
}

const (
	TestMode       = "test"
	DevMode        = "dev"
	ProductionMode = "production"
)

type Options struct {
	DatabaseURL string
	Port        int
	Mode        string
}

type ValidationError struct {
	Field     string `json:"field"`
	Validator string `json:"validator"`
}

type Response[T any] struct {
	Data T `json:"data"`
}

func Start(options Options) *API {
	ctx := context.Background()

	db, err := pgxpool.New(ctx, options.DatabaseURL)
	if err != nil {
		panic("error connecting to the database. " + err.Error())
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	api := &API{
		Context:  ctx,
		Queries:  sqlc.New(db),
		validate: validate,
		Database: db,
	}

	var r *gin.Engine
	switch options.Mode {
	case TestMode:
		gin.SetMode(gin.TestMode)
		r = gin.New()
	case DevMode:
		gin.SetMode(gin.DebugMode)
		r = gin.Default()
	default:
		gin.SetMode(gin.ReleaseMode)
		r = gin.Default()
	}

	r.GET("/health", api.getHealth)
	r.POST("/setup", api.postSetup)

	api.httpServer = &http.Server{
		Addr:    ":" + fmt.Sprint(options.Port),
		Handler: r,
	}
	go func() {
		err = api.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic("error starting server. " + err.Error())
		}
	}()
	ready := make(chan bool, 1)
	go func() {
		var res *http.Response

		for res == nil || res.StatusCode != http.StatusOK {
			res, err = http.Get("http://localhost" + api.httpServer.Addr + "/health")
			if err != nil {
				panic("error getting health check. " + err.Error())
			}
		}

		ready <- true
	}()
	<-ready

	return api
}

func (api *API) Close() {
	api.httpServer.Close()
}

func (api *API) Addr() string {
	return api.httpServer.Addr
}

func (api *API) BodyAsJSON(req any, c *gin.Context) {
	c.BindJSON(req)
	err := api.validate.Struct(req)

	if err == nil {
		return
	}

	validationErrors := err.(validator.ValidationErrors)
	apiErrors := make([]ValidationError, 0, len((validationErrors)))
	for _, validationError := range validationErrors {
		apiErrors = append(apiErrors, ValidationError{
			Field:     validationError.Field(),
			Validator: validationError.Tag(),
		})
	}

	c.JSON(http.StatusBadRequest, gin.H{"data": apiErrors})
	c.Done()
}
