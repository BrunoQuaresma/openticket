package api_test

import (
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/stretchr/testify/require"
)

func TestLogin_ValidCredentials(t *testing.T) {
	var tEnv testutil.TestEnv
	tEnv.Start()
	defer tEnv.Close()
	tEnv.Setup()
	sdk := tEnv.SDK()

	credentials := tEnv.AdminCredentials()
	r, err := sdk.Login(api.LoginRequest{
		Email:    credentials.Email,
		Password: credentials.Password,
	})

	require.NoError(t, err, "error making login request")
	require.Equal(t, 200, r.StatusCode, "unexpected status code")
	require.NotEmpty(t, r.Data.SessionToken, "session token should not be empty")
}

func TestLogin_InvalidCredentials(t *testing.T) {
	var tEnv testutil.TestEnv
	tEnv.Start()
	defer tEnv.Close()
	tEnv.Setup()
	sdk := tEnv.SDK()

	t.Run("wrong email", func(t *testing.T) {
		credentials := tEnv.AdminCredentials()
		r, err := sdk.Login(api.LoginRequest{
			Email:    "wrong" + credentials.Email,
			Password: credentials.Password,
		})

		require.NoError(t, err, "error making login request")
		require.Equal(t, 401, r.StatusCode, "unexpected status code")
	})

	t.Run("wrong password", func(t *testing.T) {
		credentials := tEnv.AdminCredentials()
		r, err := sdk.Login(api.LoginRequest{
			Email:    credentials.Email,
			Password: "wrong" + credentials.Password,
		})

		require.NoError(t, err, "error making login request")
		require.Equal(t, 401, r.StatusCode, "unexpected status code")
	})
}
