package api_test

import (
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAuthRequired(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	server := tEnv.Server()
	server.Extend(func(r *gin.Engine) {
		authorized := r.Group("/admin")
		authorized.Use(server.AuthRequired)
		authorized.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})
	})
	tEnv.Start()
	defer tEnv.Close()
	setupReq := tEnv.Setup()

	t.Run("no token", func(t *testing.T) {
		res, err := http.Get(tEnv.Server().URL() + "/admin/test")
		require.NoError(t, err, "error making admin test request")
		require.Equal(t, http.StatusUnauthorized, res.StatusCode, "expect unauthorized status code")
	})

	t.Run("invalid token", func(t *testing.T) {
		var client http.Client
		req, err := http.NewRequest("GET", tEnv.Server().URL()+"/admin/test", nil)
		require.NoError(t, err, "error creating request")
		req.Header.Set(api.SessionTokenHeader, "invalid-token")
		res, err := client.Do(req)
		require.NoError(t, err, "error making admin test request")
		require.Equal(t, http.StatusUnauthorized, res.StatusCode, "expect unauthorized status code")
	})

	t.Run("valid token", func(t *testing.T) {
		sdk := tEnv.SDK()
		var loginRes api.LoginResponse
		_, err := sdk.Login(api.LoginRequest(api.LoginRequest{
			Email:    setupReq.Email,
			Password: setupReq.Password,
		}), &loginRes)
		require.NoError(t, err, "error making login request")

		var client http.Client
		req, err := http.NewRequest("GET", tEnv.Server().URL()+"/admin/test", nil)
		require.NoError(t, err, "error creating request")
		req.Header.Set(api.SessionTokenHeader, loginRes.Data.SessionToken)
		res, err := client.Do(req)
		require.NoError(t, err, "error making admin test request")
		require.Equal(t, http.StatusOK, res.StatusCode, "expect ok status code")
	})
}

func TestLogin_ValidCredentials(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	setupReq := tEnv.Setup()
	sdk := tEnv.SDK()

	var res api.LoginResponse
	httpRes, err := sdk.Login(api.LoginRequest{
		Email:    setupReq.Email,
		Password: setupReq.Password,
	}, &res)

	require.NoError(t, err, "error making login request")
	require.Equal(t, 200, httpRes.StatusCode, "unexpected status code")
	require.NotEmpty(t, res.Data.SessionToken, "session token should not be empty")
}

func TestLogin_InvalidCredentials(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	setupReq := tEnv.Setup()
	sdk := tEnv.SDK()

	t.Run("wrong email", func(t *testing.T) {
		httpRes, err := sdk.Login(api.LoginRequest{
			Email:    "wrong" + setupReq.Email,
			Password: setupReq.Password,
		}, &api.LoginResponse{})

		require.NoError(t, err, "error making login request")
		require.Equal(t, 401, httpRes.StatusCode, "unexpected status code")
	})

	t.Run("wrong password", func(t *testing.T) {
		httpRes, err := sdk.Login(api.LoginRequest{
			Email:    setupReq.Email,
			Password: "wrong" + setupReq.Password,
		}, &api.LoginResponse{})

		require.NoError(t, err, "error making login request")
		require.Equal(t, 401, httpRes.StatusCode, "unexpected status code")
	})
}
