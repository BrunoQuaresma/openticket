package api_test

import (
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/stretchr/testify/require"
)

func TestAPI_Status(t *testing.T) {
	t.Parallel()

	t.Run("success: setup field reflects the setup state", func(t *testing.T) {
		t.Parallel()

		tEnv := testutil.NewEnv(t)
		tEnv.Start()
		sdk := tEnv.SDK()

		var res api.StatusResponse
		_, err := sdk.Status(&res)
		require.NoError(t, err, "error making status request")
		require.False(t, res.Data.Setup, "setup should be false")

		tEnv.Setup()

		sdk = tEnv.SDK()
		_, err = sdk.Status(&res)
		require.NoError(t, err, "error making status request")
		require.True(t, res.Data.Setup, "setup should be true")
	})

	t.Run("success: user field reflects the authenticated user", func(t *testing.T) {
		t.Parallel()

		tEnv := testutil.NewEnv(t)
		tEnv.Start()
		sdk := tEnv.SDK()

		var res api.StatusResponse
		_, err := sdk.Status(&res)
		require.NoError(t, err, "error making status request")
		require.Nil(t, res.Data.User, "user should be nil")

		setup := tEnv.Setup()
		sdk = tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)
		_, err = sdk.Status(&res)
		require.NoError(t, err, "error making status request")
		require.NotNil(t, res.Data.User, "user should not be nil")
	})
}
