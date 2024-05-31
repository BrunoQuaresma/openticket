package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/BrunoQuaresma/openticket/sdk"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestSetup_Validation(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()

	t.Run("required fields", func(t *testing.T) {
		var req api.SetupRequest
		var res api.Response[any]
		sdk := sdk.New(tEnv.URL())
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
		var res api.Response[any]
		sdk := sdk.New(tEnv.URL())
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

		var res api.Response[any]
		sdk := sdk.New(tEnv.URL())
		httpRes, err := sdk.Setup(req, &res)
		require.NoError(t, err, "error making request")

		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "password", "min")
	})
}

func TestSetup(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	sdk := sdk.New(tEnv.URL())

	req := api.SetupRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: testutil.FakePassword(),
	}

	var res api.Response[any]
	httpRes, err := sdk.Setup(req, &res)
	require.NoError(t, err, "error making the first request")
	require.Equal(t, http.StatusOK, httpRes.StatusCode)

	ctx := context.Background()
	firstUser, err := tEnv.DBQueries().GetUserByEmail(ctx, req.Email)
	require.NoError(t, err, "error getting the first user")
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(firstUser.PasswordHash), []byte(req.Password)), "user password should be hashed")
	require.Equal(t, database.RoleAdmin, firstUser.Role, "first user should be admin")

	req = api.SetupRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: testutil.FakePassword(),
	}
	httpRes, err = sdk.Setup(req, &res)
	require.NoError(t, err, "error making the second request")
	require.Equal(t, http.StatusNotFound, httpRes.StatusCode, "setup should return 404 if it was already done")
}
