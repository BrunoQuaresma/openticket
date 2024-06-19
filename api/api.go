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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	db         *ServerDB
	validate   *validator.Validate
	httpServer *http.Server
	router     *gin.Engine
}

const (
	TestMode       = "test"
	DevMode        = "dev"
	ProductionMode = "production"
)

type ServerOptions struct {
	DatabaseURL string
	Port        int
	Mode        string
}

type ValidationError struct {
	Field     string `json:"field"`
	Validator string `json:"validator"`
}

type Response[T any] struct {
	Data    T                 `json:"data,omitempty"`
	Errors  []ValidationError `json:"errors,omitempty"`
	Message string            `json:"message,omitempty"`
}

type ServerError struct {
	Res    Response[any]
	Status int
}

func (e ServerError) Error() string {
	return e.Res.Message
}

func NewServer(options ServerOptions) *Server {
	var server Server

	dbCtx := context.Background()
	dbConn, err := pgxpool.New(dbCtx, options.DatabaseURL)
	if err != nil {
		panic("error connecting to the database. " + err.Error())
	}
	server.db = &ServerDB{
		conn:    dbConn,
		queries: database.New(dbConn),
	}

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

	authenticated := server.router.Group("/")
	authenticated.Use(server.AuthRequired)
	{
		authenticated.POST("/users", server.createUser)
		authenticated.DELETE("/users/:id", server.deleteUser)
		authenticated.PATCH("/users/:id", server.patchUser)

		authenticated.POST("/tickets", server.createTicket)
		authenticated.GET("/tickets", server.tickets)

		authenticated.POST("/tickets/:ticketId/comments", server.createComment)
		authenticated.DELETE("/tickets/:ticketId/comments/:commentId", server.deleteComment)
		authenticated.PATCH("/tickets/:ticketId/comments/:commentId", server.patchComment)
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
			res, err = http.Get(server.URL() + "/health")
			if err != nil {
				panic("error getting health check. " + err.Error())
			}
		}

		ready <- true
	}()
	<-ready
}

func (server *Server) URL() string {
	return "http://localhost" + server.httpServer.Addr
}

func (server *Server) Close() {
	server.httpServer.Close()
	server.db.conn.Close()
}

func (server *Server) jsonReq(c *gin.Context, req any) {
	c.BindJSON(req)
	err := server.validate.Struct(req)

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

	c.AbortWithStatusJSON(http.StatusBadRequest, Response[any]{Errors: apiErrors})
}

type ServerDB struct {
	conn    *pgxpool.Pool
	queries *database.Queries
}

type txFn func(ctx context.Context, qtx *database.Queries, tx pgx.Tx) error

func (db *ServerDB) tx(fn txFn) error {
	ctx := context.Background()
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = fn(ctx, db.queries.WithTx(tx), tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

type PermissionDeniedError struct {
	Message string
}

func (e PermissionDeniedError) Error() string {
	return "permission denied: " + e.Message
}
