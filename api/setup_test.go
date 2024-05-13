package api_test

import (
	"testing"

	"github.com/BrunoQuaresma/openticket/api/test"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	s := test.TestServer{}
	s.Start()
	defer s.Close()
	require.True(t, true)
}
