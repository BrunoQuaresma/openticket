package api_test

import (
	"testing"

	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	tEnv := testutil.TestEnv{}
	tEnv.Start()
	defer tEnv.Close()
	require.True(t, true)
}
