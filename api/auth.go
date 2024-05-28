package api

import (
	"context"
	"time"

	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse = Response[struct {
	SessionToken string `json:"session_token"`
}]

func (server *Server) login(c *gin.Context) {
	var req LoginRequest
	server.ParseJSONRequest(c, &req)

	ctx := context.Background()
	dbQueries := server.DBQueries()
	user, err := dbQueries.GetUserByEmail(ctx, req.Email)

	if err != nil {
		if err == pgx.ErrNoRows {
			c.AbortWithStatusJSON(401, gin.H{"message": "invalid email or password"})
			return
		}

		c.AbortWithError(500, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "invalid email or password"})
		return
	}

	t := uuid.NewString()
	tokenHash, err := bcrypt.GenerateFromPassword([]byte(t), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	_, err = dbQueries.CreateSession(ctx, database.CreateSessionParams{
		UserID:    user.ID,
		TokenHash: string(tokenHash),
		ExpiresAt: pgtype.Timestamp{
			Time:  time.Now().AddDate(0, 0, 30).UTC(),
			Valid: true,
		},
	})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	var res LoginResponse
	res.Data.SessionToken = t

	c.JSON(200, res)
}
