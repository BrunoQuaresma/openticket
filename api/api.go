package api

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/BrunoQuaresma/openticket/api/database"
)

type Server struct {
	db         *database.Connection
	validate   *validator.Validate
	httpServer *http.Server
	router     *gin.Engine
}

const (
	TestMode       = "test"
	DevMode        = "dev"
	ProductionMode = "production"
)

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

func NewServer(port int, database *database.Connection, mode string) *Server {
	server := Server{db: database}

	server.validate = validator.New(validator.WithRequiredStructEnabled())
	server.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	switch mode {
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
		authenticated.GET("/tickets/:ticketId", server.ticket)
		authenticated.DELETE("/tickets/:ticketId", server.deleteTicket)
		authenticated.PATCH("/tickets/:ticketId", server.patchTicket)
		authenticated.PATCH("/tickets/:ticketId/status", server.patchTicketStatus)

		authenticated.POST("/tickets/:ticketId/comments", server.createComment)
		authenticated.DELETE("/tickets/:ticketId/comments/:commentId", server.deleteComment)
		authenticated.PATCH("/tickets/:ticketId/comments/:commentId", server.patchComment)

		authenticated.POST("/tickets/:ticketId/assignments", server.createAssignment)
		authenticated.DELETE("/tickets/:ticketId/assignments/:assignmentId", server.deleteAssignment)
	}

	server.httpServer = &http.Server{
		Addr:    ":" + fmt.Sprint(port),
		Handler: server.router,
	}

	return &server
}

func (server *Server) Extend(f func(r *gin.Engine)) {
	f(server.router)
}

func (server *Server) Start() {
	err := server.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("error starting server. " + err.Error())
	}
}

func (server *Server) URL() string {
	return "http://localhost" + server.httpServer.Addr
}

func (server *Server) Close() {
	server.httpServer.Close()
	server.db.Close()
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

type PermissionDeniedError struct {
	Message string
}

func (e PermissionDeniedError) Error() string {
	return "permission denied: " + e.Message
}
