package server_test

import (
	"testing"

	"github.com/BrunoQuaresma/openticket/testutils"
	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	c := testutils.NewTestDatabaseConfig()
	c.Database = "setup-db"
	stop := testutils.RunTestServer(c)
	defer stop()
	require.True(t, true)
}
