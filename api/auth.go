package api

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

const SessionTokenHeader = "OPENTICKET-SESSION-TOKEN"
const userCtxKey = "user"

func (server *Server) AuthRequired(c *gin.Context) {
	sessionToken := c.Request.Header.Get(SessionTokenHeader)
	queries := server.DBQueries()
	ctx := context.Background()
	sum := sha256.Sum256([]byte(sessionToken))
	tokenHash := base64.URLEncoding.EncodeToString(sum[:])
	session, err := queries.GetSessionByTokenHash(ctx, tokenHash)

	if err != nil {
		if err == pgx.ErrNoRows {
			c.AbortWithStatus(401)
		} else {
			c.AbortWithError(500, err)
		}
		return
	}
	if session.ExpiresAt.Time.Before(time.Now()) {
		c.AbortWithStatus(401)
		return
	}

	user, err := queries.GetUserByID(ctx, session.UserID)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Set(userCtxKey, user)
	c.Next()
}

func (server *Server) AuthUser(c *gin.Context) database.User {
	user, err := c.Get(userCtxKey)
	if !err {
		c.AbortWithError(http.StatusInternalServerError, errors.New("user not found in context"))
		return database.User{}
	}
	return user.(database.User)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse = Response[struct {
	SessionToken string `json:"session_token"`
}]

func (server *Server) login(c *gin.Context) {
	var req LoginRequest
	server.JSONRequest(c, &req)

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

	token, err := secureToken()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	sum := sha256.Sum256([]byte(token))
	tokenHash := base64.URLEncoding.EncodeToString(sum[:])
	_, err = dbQueries.CreateSession(ctx, database.CreateSessionParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
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
	res.Data.SessionToken = token

	c.JSON(200, res)
}

func secureToken() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
