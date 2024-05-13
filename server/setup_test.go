package server_test

import (
	"testing"

	"github.com/BrunoQuaresma/openticket/testutils"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	s := testutils.TestServer{}
	s.Start()
	defer s.Close()
	require.True(t, true)
}
