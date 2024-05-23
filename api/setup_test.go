package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestSetup_Validation(t *testing.T) {
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
		req := api.PostSetupRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    "invalid-email",
			Password: fakePassword(),
		}
		var apiRes api.Response[[]api.ValidationError]
		r, err := postSetup(&tEnv, req, &apiRes)
		require.NoError(t, err, "error making request")

		require.Equal(t, http.StatusBadRequest, r.StatusCode)
		require.True(t, testutil.HasValidationError(apiRes.Data, "email", "email"))
	})

	t.Run("valid password", func(t *testing.T) {
		req := api.PostSetupRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: "no8char",
		}
		var apiRes api.Response[[]api.ValidationError]
		r, err := postSetup(&tEnv, req, &apiRes)
		require.NoError(t, err, "error making request")

		require.Equal(t, http.StatusBadRequest, r.StatusCode)
		require.True(t, testutil.HasValidationError(apiRes.Data, "password", "min"))
	})
}

func TestSetup(t *testing.T) {
	tEnv := testutil.TestEnv{}
	tEnv.Start()
	defer tEnv.Close()

	req := api.PostSetupRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: fakePassword(),
	}
	r, err := postSetup(&tEnv, req, nil)
	require.NoError(t, err, "error making the first request")
	require.Equal(t, http.StatusCreated, r.StatusCode)

	ctx := context.Background()
	firstUser, err := tEnv.Server.Queries.GetUserByEmail(ctx, req.Email)
	require.NoError(t, err, "error getting the first user")
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(firstUser.Hash), []byte(req.Password)), "user password should be hashed")
	require.Equal(t, database.RoleAdmin, firstUser.Role, "first user should be admin")

	req = api.PostSetupRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: fakePassword(),
	}
	r, err = postSetup(&tEnv, req, nil)
	require.NoError(t, err, "error making the second request")
	require.Equal(t, http.StatusNotFound, r.StatusCode, "setup should return 404 if it was already done")
}

func postSetup(tEnv *testutil.TestEnv, req api.PostSetupRequest, res any) (*http.Response, error) {
	b, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	r, err := http.Post(tEnv.URL()+"/setup", "application/json", bytes.NewBuffer(b))
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
