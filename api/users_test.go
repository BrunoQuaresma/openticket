package api_test

import (
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/BrunoQuaresma/openticket/sdk"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Authentication(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	tEnv.Setup()
	sdk := sdk.New(tEnv.URL())

	httpRes, err := sdk.CreateUser(api.CreateUserRequest{}, &api.CreateUserResponse{})
	require.NoError(t, err, "error making create user request")
	require.Equal(t, http.StatusUnauthorized, httpRes.StatusCode)
}

func TestCreateUser_Validation(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	setupReq := tEnv.Setup()
	sdk := sdk.New(tEnv.URL())
	var loginRes api.LoginResponse
	_, err := sdk.Login(api.LoginRequest(api.LoginRequest{
		Email:    setupReq.Email,
		Password: setupReq.Password,
	}), &loginRes)
	require.NoError(t, err, "error making login request")
	sdk.Authenticate(loginRes.Data.SessionToken)

	t.Run("required fields", func(t *testing.T) {
		var req api.CreateUserRequest
		var res api.CreateUserResponse
		httpRes, err := sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "name", "required")
		testutil.RequireValidationError(t, res.Errors, "username", "required")
		testutil.RequireValidationError(t, res.Errors, "email", "required")
		testutil.RequireValidationError(t, res.Errors, "password", "required")
		testutil.RequireValidationError(t, res.Errors, "role", "required")
	})

	t.Run("valid email", func(t *testing.T) {
		req := api.CreateUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    "invalid-email",
			Password: testutil.FakePassword(),
			Role:     "member",
		}
		var res api.CreateUserResponse
		httpRes, err := sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "email", "email")
	})

	t.Run("valid password", func(t *testing.T) {
		req := api.CreateUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: "no8char",
			Role:     "member",
		}
		var res api.CreateUserResponse
		httpRes, err := sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "password", "min")
	})

	t.Run("valid role", func(t *testing.T) {
		req := api.CreateUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: testutil.FakePassword(),
			Role:     "invalid-role",
		}

		var res api.CreateUserResponse
		httpRes, err := sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "role", "oneof")
	})

	t.Run("unique email", func(t *testing.T) {
		req := api.CreateUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: testutil.FakePassword(),
			Role:     "member",
		}

		var res api.CreateUserResponse
		httpRes, err := sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusCreated, httpRes.StatusCode)

		// Use a different username to avoid unique constraint violation. We only
		// care about email.
		req.Username = gofakeit.Username()
		httpRes, err = sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "email", "unique")
	})

	t.Run("unique username", func(t *testing.T) {
		req := api.CreateUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: testutil.FakePassword(),
			Role:     "member",
		}

		var res api.CreateUserResponse
		httpRes, err := sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusCreated, httpRes.StatusCode)

		// Use a different email to avoid unique constraint violation. We only care
		// about username.
		req.Email = gofakeit.Email()
		httpRes, err = sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "username", "unique")
	})
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	setupReq := tEnv.Setup()
	sdk := sdk.New(tEnv.URL())
	var loginRes api.LoginResponse
	_, err := sdk.Login(api.LoginRequest(api.LoginRequest{
		Email:    setupReq.Email,
		Password: setupReq.Password,
	}), &loginRes)
	require.NoError(t, err, "error making login request")
	sdk.Authenticate(loginRes.Data.SessionToken)

	req := api.CreateUserRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: testutil.FakePassword(),
		Role:     "member",
	}

	var res api.CreateUserResponse
	httpRes, err := sdk.CreateUser(req, &res)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	require.NotEmpty(t, res.Data.ID)
	require.Equal(t, req.Name, res.Data.Name)
	require.Equal(t, req.Username, res.Data.Username)
	require.Equal(t, req.Email, res.Data.Email)
	require.Equal(t, req.Role, res.Data.Role)
}
