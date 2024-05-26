package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestSetup_Validation(t *testing.T) {
	var tEnv testutil.TestEnv
	tEnv.Start()
	defer tEnv.Close()

	t.Run("required fields", func(t *testing.T) {
		var req api.SetupRequest
		r, err := tEnv.SDK().Setup(req)
		require.NoError(t, err, "error making request")

		require.Equal(t, http.StatusBadRequest, r.StatusCode)
		require.True(t, testutil.HasValidationError(r, "name", "required"))
		require.True(t, testutil.HasValidationError(r, "username", "required"))
		require.True(t, testutil.HasValidationError(r, "email", "required"))
		require.True(t, testutil.HasValidationError(r, "password", "required"))
	})

	t.Run("valid email", func(t *testing.T) {
		req := api.SetupRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    "invalid-email",
			Password: testutil.FakePassword(),
		}
		r, err := tEnv.SDK().Setup(req)
		require.NoError(t, err, "error making request")

		require.Equal(t, http.StatusBadRequest, r.StatusCode)
		require.True(t, testutil.HasValidationError(r, "email", "email"))
	})

	t.Run("valid password", func(t *testing.T) {
		req := api.SetupRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: "no8char",
		}

		r, err := tEnv.SDK().Setup(req)
		require.NoError(t, err, "error making request")

		require.Equal(t, http.StatusBadRequest, r.StatusCode)
		require.True(t, testutil.HasValidationError(r, "password", "min"))
	})
}

func TestSetup(t *testing.T) {
	var tEnv testutil.TestEnv
	tEnv.Start()
	defer tEnv.Close()
	sdk := tEnv.SDK()

	req := api.SetupRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: testutil.FakePassword(),
	}
	r, err := sdk.Setup(req)
	require.NoError(t, err, "error making the first request")
	require.Equal(t, http.StatusOK, r.StatusCode)

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
	r, err = sdk.Setup(req)
	require.NoError(t, err, "error making the second request")
	require.Equal(t, http.StatusNotFound, r.StatusCode, "setup should return 404 if it was already done")
}
