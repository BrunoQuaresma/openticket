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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	db         *pgxpool.Pool
	dbQueries  *sqlc.Queries
	validate   *validator.Validate
	httpServer *http.Server
	router     *gin.Engine
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
	Data   T                 `json:"data,omitempty"`
	Errors []ValidationError `json:"errors"`
}

func New(options Options) *Server {
	var server Server

	dbCtx := context.Background()
	db, err := pgxpool.New(dbCtx, options.DatabaseURL)
	if err != nil {
		panic("error connecting to the database. " + err.Error())
	}
	server.db = db
	server.dbQueries = sqlc.New(db)

	server.validate = validator.New(validator.WithRequiredStructEnabled())
	server.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	switch options.Mode {
	case TestMode:
		gin.SetMode(gin.TestMode)
		server.router = gin.New()
	case DevMode:
		gin.SetMode(gin.DebugMode)
		server.router = gin.Default()
	default:
		gin.SetMode(gin.ReleaseMode)
		server.router = gin.Default()
	}

	server.router.GET("/health", server.health)
	server.router.POST("/setup", server.setup)
	server.router.POST("/login", server.login)
	server.router.POST("/tickets", server.createTicket)

	authenticated := server.router.Group("/")
	{
		authenticated.Use(server.AuthRequired)
		authenticated.POST("/users", server.createUser)
	}

	server.httpServer = &http.Server{
		Addr:    ":" + fmt.Sprint(options.Port),
		Handler: server.router,
	}

	return &server
}

func (server *Server) Extend(f func(r *gin.Engine)) {
	f(server.router)
}

func (server *Server) Start() {
	go func() {
		defer server.Close()
		err := server.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic("error starting server. " + err.Error())
		}
	}()
	ready := make(chan bool, 1)
	go func() {
		var (
			res *http.Response
			err error
		)

		for res == nil || res.StatusCode != http.StatusOK {
			res, err = http.Get("http://localhost" + server.Addr() + "/health")
			if err != nil {
				panic("error getting health check. " + err.Error())
			}
		}

		ready <- true
	}()
	<-ready
}

func (api *Server) Close() {
	api.httpServer.Close()
	api.db.Close()
}

func (api *Server) Addr() string {
	return api.httpServer.Addr
}

func (api *Server) BeginTX(ctx context.Context) (pgx.Tx, error) {
	return api.db.Begin(ctx)
}

func (api *Server) DBQueries() *sqlc.Queries {
	return api.dbQueries
}

func (api *Server) ParseJSONRequest(c *gin.Context, req any) {
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

	c.JSON(http.StatusBadRequest, Response[any]{Errors: apiErrors})
	c.Done()
}
