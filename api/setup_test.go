package api_test

import (
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestAPI_Setup(t *testing.T) {
	t.Parallel()

	t.Run("validation", func(t *testing.T) {
		t.Parallel()

		tEnv := testutil.NewEnv(t)
		tEnv.Start()
		t.Cleanup(tEnv.Close)

		t.Run("error: no required fields", func(t *testing.T) {
			t.Parallel()

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

		t.Run("error: invalid email", func(t *testing.T) {
			t.Parallel()

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

		t.Run("error: invalid password", func(t *testing.T) {
			t.Parallel()

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
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		tEnv := testutil.NewEnv(t)
		tEnv.Start()
		t.Cleanup(tEnv.Close)
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
	})

	t.Run("error: setup can't run twice", func(t *testing.T) {
		t.Parallel()

		tEnv := testutil.NewEnv(t)
		tEnv.Start()
		t.Cleanup(tEnv.Close)
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
	})
}
