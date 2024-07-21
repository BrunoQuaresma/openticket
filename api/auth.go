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

	sqlc "github.com/BrunoQuaresma/openticket/api/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

const TokenHeader = "OPENTICKET-TOKEN"
const TokenCookie = "openticket-token"
const userCtxKey = "user"

func (server *Server) AuthRequired(c *gin.Context) {
	user, err := server.AuthUser(c)

	if user == nil {
		c.AbortWithStatus(401)
		return
	}

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Set(userCtxKey, user)
	c.Next()
}

func (server *Server) AuthUserFromContext(c *gin.Context) sqlc.User {
	user, err := c.Get(userCtxKey)
	if !err {
		c.AbortWithError(http.StatusInternalServerError, errors.New("user not found in context. Ensure the AuthRequired middleware is used in the route calling this function"))
		return sqlc.User{}
	}
	return user.(sqlc.User)
}

func (server *Server) AuthUser(c *gin.Context) (user *sqlc.User, err error) {
	sessionToken := c.Request.Header.Get(TokenHeader)
	if sessionToken == "" {
		sessionToken, err = c.Cookie(TokenCookie)
		if err != nil {
			return nil, nil
		}
	}
	sum := sha256.Sum256([]byte(sessionToken))
	tokenHash := base64.URLEncoding.EncodeToString(sum[:])
	session, err := server.db.Queries().GetSessionByTokenHash(c, tokenHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if session.ExpiresAt.Time.Before(time.Now()) {
		return nil, nil
	}

	result, err := server.db.Queries().GetUserByID(c, session.UserID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginData struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type LoginResponse = Response[LoginData]

func (server *Server) login(c *gin.Context) {
	var req LoginRequest
	server.jsonReq(c, &req)

	ctx := context.Background()
	user, err := server.db.Queries().GetUserByEmail(ctx, req.Email)

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
	maxAge := time.Hour * 24 * 30
	tokenExpiration := time.Now().Add(maxAge)
	_, err = server.db.Queries().CreateSession(ctx, sqlc.CreateSessionParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamp{
			Time:  tokenExpiration.UTC(),
			Valid: true,
		},
	})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.SetCookie(TokenCookie, token, int(maxAge), "/", "", false, true)
	c.JSON(200, LoginResponse{
		Data: LoginData{
			User: User{
				ID:       user.ID,
				Name:     user.Name,
				Username: user.Username,
				Email:    user.Email,
				Role:     string(user.Role),
			},
			Token: token,
		},
	})
}

func secureToken() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
