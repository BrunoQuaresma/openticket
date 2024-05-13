package api_test

import (
	"testing"

	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	s := testutil.TestServer{}
	s.Start()
	defer s.Close()
	require.True(t, true)
}
