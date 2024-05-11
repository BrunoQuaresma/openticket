package server_test

import (
	"testing"

	"github.com/BrunoQuaresma/openticket/testutils"
	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	dbConf := testutils.NewTestDatabaseConfig()
	dbConf.Database = "setup-db"
	stop := testutils.RunTestServer(testutils.TestServerConfig{
		Database: dbConf,
	})
	defer stop()
	require.True(t, true)
}
