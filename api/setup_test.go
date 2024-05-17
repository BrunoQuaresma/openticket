package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestSetupValidation(t *testing.T) {
	tEnv := testutil.TestEnv{}
	tEnv.Start()
	defer tEnv.Close()

	t.Run("required fields", func(t *testing.T) {
		var body api.PostSetupRequest
		var apiRes api.Response[[]api.ValidationError]
		r, err := postSetup(&tEnv, body, &apiRes)
		require.NoError(t, err, "error making request")

		require.Equal(t, http.StatusBadRequest, r.StatusCode)
		require.True(t, testutil.HasValidationError(apiRes.Data, "name", "required"))
		require.True(t, testutil.HasValidationError(apiRes.Data, "username", "required"))
		require.True(t, testutil.HasValidationError(apiRes.Data, "email", "required"))
		require.True(t, testutil.HasValidationError(apiRes.Data, "password", "required"))
	})

	t.Run("valid email", func(t *testing.T) {
		body := api.PostSetupRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    "invalid-email",
			Password: fakePassword(),
		}
		var apiRes api.Response[[]api.ValidationError]
		r, err := postSetup(&tEnv, body, &apiRes)
		require.NoError(t, err, "error making request")

		require.Equal(t, http.StatusBadRequest, r.StatusCode)
		require.True(t, testutil.HasValidationError(apiRes.Data, "email", "email"))
	})

	t.Run("valid password", func(t *testing.T) {
		body := api.PostSetupRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: "no8char",
		}
		var apiRes api.Response[[]api.ValidationError]
		r, err := postSetup(&tEnv, body, &apiRes)
		require.NoError(t, err, "error making request")

		require.Equal(t, http.StatusBadRequest, r.StatusCode)
		require.True(t, testutil.HasValidationError(apiRes.Data, "password", "min"))
	})
}

func postSetup(tEnv *testutil.TestEnv, body api.PostSetupRequest, res any) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	r, err := http.Post(tEnv.URL+"/setup", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(res)
	return r, nil
}

func fakePassword() string {
	return gofakeit.Password(true, true, true, true, false, 15)
}
