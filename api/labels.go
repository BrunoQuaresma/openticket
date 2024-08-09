package api

import (
	"net/http"

	"github.com/BrunoQuaresma/openticket/api/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Label struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LabelsResponse = Response[[]Label]

func (s *Server) labels(c *gin.Context) {
	labels, err := s.db.Queries().GetLabels(c)
	if err != nil && err != pgx.ErrNoRows {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response[any]{Message: "failed to get labels"})
		return
	}

	var data []Label
	for _, label := range labels {
		data = append(data, Label{
			ID:   int(label.ID),
			Name: label.Name,
		})
	}
	c.JSON(http.StatusOK, LabelsResponse{
		Data: data,
	})
}

type CreateLabelRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateLabelResponse = Response[Label]

func (s *Server) createLabel(c *gin.Context) {
	user := s.AuthUserFromContext(c)

	var req CreateLabelRequest
	s.jsonReq(c, &req)

	label, err := s.db.Queries().CreateLabel(c, sqlc.CreateLabelParams{
		Name:      req.Name,
		CreatedBy: user.ID,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response[any]{Message: "failed to create label"})
		return
	}

	c.JSON(http.StatusCreated, CreateLabelResponse{
		Data: Label{
			ID:   int(label.ID),
			Name: label.Name,
		},
	})
}
