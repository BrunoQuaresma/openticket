package api_test

import (
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAPI_AuthRequired(t *testing.T) {
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
	setup := tEnv.Setup()

	t.Run("unauthorized: no token", func(t *testing.T) {
		t.Parallel()

		res, err := http.Get(tEnv.Server().URL() + "/admin/test")
		require.NoError(t, err, "error making admin test request")
		require.Equal(t, http.StatusUnauthorized, res.StatusCode, "expect unauthorized status code")
	})

	t.Run("unauthorized: invalid token", func(t *testing.T) {
		t.Parallel()

		var client http.Client
		req, err := http.NewRequest("GET", tEnv.Server().URL()+"/admin/test", nil)
		require.NoError(t, err, "error creating request")
		req.Header.Set(api.TokenHeader, "invalid-token")

		res, err := client.Do(req)
		require.NoError(t, err, "error making admin test request")
		require.Equal(t, http.StatusUnauthorized, res.StatusCode, "expect unauthorized status code")
	})

	t.Run("unauthorized: invalid cookie", func(t *testing.T) {
		t.Parallel()

		var client http.Client
		req, err := http.NewRequest("GET", tEnv.Server().URL()+"/admin/test", nil)
		require.NoError(t, err, "error creating request")
		req.Header.Set(api.TokenHeader, "invalid-token")
		req.AddCookie(&http.Cookie{Name: api.TokenCookie, Value: "invalid-token"})

		res, err := client.Do(req)
		require.NoError(t, err, "error making admin test request")
		require.Equal(t, http.StatusUnauthorized, res.StatusCode, "expect unauthorized status code")
	})

	t.Run("authorized: valid token", func(t *testing.T) {
		t.Parallel()

		sdk := tEnv.SDK()
		var loginRes api.LoginResponse
		_, err := sdk.Login(api.LoginRequest(api.LoginRequest{
			Email:    setup.Req().Email,
			Password: setup.Req().Password,
		}), &loginRes)
		require.NoError(t, err, "error making login request")

		var client http.Client
		req, err := http.NewRequest("GET", tEnv.Server().URL()+"/admin/test", nil)
		require.NoError(t, err, "error creating request")
		req.Header.Set(api.TokenHeader, loginRes.Data.Token)
		res, err := client.Do(req)
		require.NoError(t, err, "error making admin test request")
		require.Equal(t, http.StatusOK, res.StatusCode, "expect ok status code")
	})

	t.Run("authorized: valid cookie", func(t *testing.T) {
		t.Parallel()

		sdk := tEnv.SDK()
		var loginRes api.LoginResponse
		_, err := sdk.Login(api.LoginRequest(api.LoginRequest{
			Email:    setup.Req().Email,
			Password: setup.Req().Password,
		}), &loginRes)
		require.NoError(t, err, "error making login request")

		var client http.Client
		req, err := http.NewRequest("GET", tEnv.Server().URL()+"/admin/test", nil)
		require.NoError(t, err, "error creating request")
		req.AddCookie(&http.Cookie{Name: api.TokenCookie, Value: loginRes.Data.Token})
		res, err := client.Do(req)
		require.NoError(t, err, "error making admin test request")
		require.Equal(t, http.StatusOK, res.StatusCode, "expect ok status code")
	})
}

func TestAPI_Login(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	setup := tEnv.Setup()
	sdk := tEnv.SDK()

	t.Run("success: valid credentials", func(t *testing.T) {
		t.Parallel()

		var res api.LoginResponse
		httpRes, err := sdk.Login(api.LoginRequest{
			Email:    setup.Req().Email,
			Password: setup.Req().Password,
		}, &res)

		require.NoError(t, err, "error making login request")
		require.Equal(t, 200, httpRes.StatusCode, "unexpected status code")
		require.NotEmpty(t, res.Data.Token, "session token should not be empty")

		var tokenCookie *http.Cookie
		for _, cookie := range httpRes.Cookies() {
			if cookie.Name == api.TokenCookie {
				tokenCookie = cookie
				break
			}
		}
		require.NotNil(t, tokenCookie, "token cookie not found")
	})

	t.Run("error: wrong email", func(t *testing.T) {
		t.Parallel()

		httpRes, err := sdk.Login(api.LoginRequest{
			Email:    "wrong" + setup.Req().Email,
			Password: setup.Req().Password,
		}, &api.LoginResponse{})

		require.NoError(t, err, "error making login request")
		require.Equal(t, 401, httpRes.StatusCode, "unexpected status code")
	})

	t.Run("error: wrong password", func(t *testing.T) {
		t.Parallel()

		httpRes, err := sdk.Login(api.LoginRequest{
			Email:    setup.Req().Email,
			Password: "wrong" + setup.Req().Password,
		}, &api.LoginResponse{})

		require.NoError(t, err, "error making login request")
		require.Equal(t, 401, httpRes.StatusCode, "unexpected status code")
	})
}
