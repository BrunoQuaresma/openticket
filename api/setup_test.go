package api_test

import (
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestSetup_Validation(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()

	t.Run("required fields", func(t *testing.T) {
		var (
			req api.SetupRequest
			res api.SetupResponse
		)
		sdk := tEnv.SDK()
		httpRes, err := sdk.Setup(req, &res)
		require.NoError(t, err, "error making request")

		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "name", "required")
		testutil.RequireValidationError(t, res.Errors, "username", "required")
		testutil.RequireValidationError(t, res.Errors, "email", "required")
		testutil.RequireValidationError(t, res.Errors, "password", "required")
	})

	t.Run("valid email", func(t *testing.T) {
		req := api.SetupRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    "invalid-email",
			Password: testutil.FakePassword(),
		}
		var res api.SetupResponse
		sdk := tEnv.SDK()
		httpRes, err := sdk.Setup(req, &res)
		require.NoError(t, err, "error making request")

		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "email", "email")
	})

	t.Run("valid password", func(t *testing.T) {
		req := api.SetupRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: "no8char",
		}

		var res api.SetupResponse
		sdk := tEnv.SDK()
		httpRes, err := sdk.Setup(req, &res)
		require.NoError(t, err, "error making request")

		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "password", "min")
	})
}

func TestSetup_Success(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	sdk := tEnv.SDK()

	req := api.SetupRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: testutil.FakePassword(),
	}

	var res api.SetupResponse
	httpRes, err := sdk.Setup(req, &res)
	require.NoError(t, err, "error making the first request")
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.NotEmpty(t, res.Data.ID)
	require.Equal(t, "admin", res.Data.Role)
}

func TestSetup_CantRunTwice(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	tEnv.Setup()
	sdk := tEnv.SDK()

	req := api.SetupRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: testutil.FakePassword(),
	}

	httpRes, err := sdk.Setup(req, nil)
	require.NoError(t, err, "error making the first request")
	require.Equal(t, http.StatusForbidden, httpRes.StatusCode)
}

func TestSetup_Login(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	setup := tEnv.Setup()
	sdk := tEnv.SDK()

	var res api.LoginResponse
	httpRes, err := sdk.Login(api.LoginRequest{
		Email:    setup.Req().Email,
		Password: setup.Req().Password,
	}, &res)

	require.NoError(t, err, "error making login request")
	require.Equal(t, 200, httpRes.StatusCode, "unexpected status code")
}
