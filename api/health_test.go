package api_test

import (
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/stretchr/testify/require"
)

func TestAPI_Health(t *testing.T) {
	t.Parallel()

	t.Run("setup field reflects the setup state", func(t *testing.T) {
		t.Parallel()

		tEnv := testutil.NewEnv(t)
		tEnv.Start()
		sdk := tEnv.SDK()

		var res api.HealthResponse
		_, err := sdk.Health(&res)
		require.NoError(t, err, "error making health request")
		require.False(t, res.Data.Setup, "setup should be false")

		tEnv.Setup()

		sdk = tEnv.SDK()
		_, err = sdk.Health(&res)
		require.NoError(t, err, "error making health request")
		require.True(t, res.Data.Setup, "setup should be true")
	})
}
