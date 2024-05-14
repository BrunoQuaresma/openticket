package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	tEnv := testutil.TestEnv{}
	tEnv.Start()
	defer tEnv.Close()

	body := []byte(`{}`)

	r, err := http.Post(tEnv.URL+"/setup", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err, "error making request")

	defer r.Body.Close()
	var apiRes api.Response[[]api.ValidationError]
	json.NewDecoder(r.Body).Decode(&apiRes)

	require.Equal(t, http.StatusBadRequest, r.StatusCode)
	require.True(t, testutil.HasValidationError(apiRes.Data, "name", "required"))
	require.True(t, testutil.HasValidationError(apiRes.Data, "username", "required"))
	require.True(t, testutil.HasValidationError(apiRes.Data, "email", "required"))
	require.True(t, testutil.HasValidationError(apiRes.Data, "password", "required"))
}
