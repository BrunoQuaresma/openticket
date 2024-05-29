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

	tEnv := testutil.NewEnv()
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
	tEnv.Setup()

	t.Run("no token", func(t *testing.T) {
		res, err := http.Get(tEnv.URL() + "/admin/test")
		require.NoError(t, err, "error making admin test request")
		require.Equal(t, http.StatusUnauthorized, res.StatusCode, "expect unauthorized status code")
	})

	t.Run("invalid token", func(t *testing.T) {
		var client http.Client
		req, err := http.NewRequest("GET", tEnv.URL()+"/admin/test", nil)
		require.NoError(t, err, "error creating request")
		req.Header.Set(api.SessionTokenHeader, "invalid-token")
		res, err := client.Do(req)
		require.NoError(t, err, "error making admin test request")
		require.Equal(t, http.StatusUnauthorized, res.StatusCode, "expect unauthorized status code")
	})

	t.Run("valid token", func(t *testing.T) {
		sdk := tEnv.SDK()
		credentials := tEnv.AdminCredentials()
		var loginRes api.LoginResponse
		_, err := sdk.Login(api.LoginRequest(credentials), &loginRes)
		require.NoError(t, err, "error making login request")

		var client http.Client
		req, err := http.NewRequest("GET", tEnv.URL()+"/admin/test", nil)
		require.NoError(t, err, "error creating request")
		req.Header.Set(api.SessionTokenHeader, loginRes.Data.SessionToken)
		res, err := client.Do(req)
		require.NoError(t, err, "error making admin test request")
		require.Equal(t, http.StatusOK, res.StatusCode, "expect ok status code")
	})
}

func TestLogin_ValidCredentials(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv()
	tEnv.Start()
	defer tEnv.Close()
	tEnv.Setup()
	sdk := tEnv.SDK()

	credentials := tEnv.AdminCredentials()
	var res api.LoginResponse
	httpRes, err := sdk.Login(api.LoginRequest{
		Email:    credentials.Email,
		Password: credentials.Password,
	}, &res)

	require.NoError(t, err, "error making login request")
	require.Equal(t, 200, httpRes.StatusCode, "unexpected status code")
	require.NotEmpty(t, res.Data.SessionToken, "session token should not be empty")
}

func TestLogin_InvalidCredentials(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv()
	tEnv.Start()
	defer tEnv.Close()
	tEnv.Setup()
	sdk := tEnv.SDK()

	t.Run("wrong email", func(t *testing.T) {
		credentials := tEnv.AdminCredentials()
		httpRes, err := sdk.Login(api.LoginRequest{
			Email:    "wrong" + credentials.Email,
			Password: credentials.Password,
		}, &api.LoginResponse{})

		require.NoError(t, err, "error making login request")
		require.Equal(t, 401, httpRes.StatusCode, "unexpected status code")
	})

	t.Run("wrong password", func(t *testing.T) {
		credentials := tEnv.AdminCredentials()
		httpRes, err := sdk.Login(api.LoginRequest{
			Email:    credentials.Email,
			Password: "wrong" + credentials.Password,
		}, &api.LoginResponse{})

		require.NoError(t, err, "error making login request")
		require.Equal(t, 401, httpRes.StatusCode, "unexpected status code")
	})
}
