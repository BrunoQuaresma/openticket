package server_test

import (
	"testing"

	"github.com/BrunoQuaresma/openticket/testutils"
	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	stop := testutils.RunTestServer(testutils.TestServerConfig{})
	defer stop()
	require.True(t, true)
}
