package api_test

import (
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/stretchr/testify/require"
)

func TestHealth_SetupDone(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	tEnv.Setup()
	sdk := tEnv.SDK()

	var res api.HealthResponse
	_, err := sdk.Health(&res)
	require.NoError(t, err, "error making health request")
	require.True(t, res.Data.Setup, "health check failed")
}

func TestHealth_NoSetup(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	sdk := tEnv.SDK()

	var res api.HealthResponse
	_, err := sdk.Health(&res)
	require.NoError(t, err, "error making health request")
	require.False(t, res.Data.Setup, "health check failed")
}
