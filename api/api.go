package api

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type API struct {
	Context  context.Context
	Queries  *database.Queries
	validate *validator.Validate
}

type Options struct {
	DatabaseURL string
	Debug       bool
	Port        int
}

type ValidationError struct {
	Field     string `json:"field"`
	Validator string `json:"validator"`
}

type Response[T any] struct {
	Data T `json:"data"`
}

func Start(options Options) *http.Server {
	ctx := context.Background()

	d, err := pgxpool.New(ctx, options.DatabaseURL)
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

	api := API{
		Context:  ctx,
		Queries:  database.New(d),
		validate: validate,
	}

	if !options.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.GET("/health", api.getHealth)
	r.POST("/setup", api.postSetup)

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

func (api *API) BodyAsJSON(body any, c *gin.Context) {
	c.BindJSON(body)
	err := api.validate.Struct(body)

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
