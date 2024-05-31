package api_test

import (
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Authentication(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	tEnv.Setup()
	sdk := tEnv.SDK()

	httpRes, err := sdk.CreateUser(api.CreateUserRequest{}, &api.CreateUserResponse{})
	require.NoError(t, err, "error making create user request")
	require.Equal(t, http.StatusUnauthorized, httpRes.StatusCode)
}

func TestCreateUser_Validation(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	t.Run("required fields", func(t *testing.T) {
		var req api.CreateUserRequest
		var res api.CreateUserResponse
		httpRes, err := sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "name", "required")
		testutil.RequireValidationError(t, res.Errors, "username", "required")
		testutil.RequireValidationError(t, res.Errors, "email", "required")
		testutil.RequireValidationError(t, res.Errors, "password", "required")
		testutil.RequireValidationError(t, res.Errors, "role", "required")
	})

	t.Run("valid email", func(t *testing.T) {
		req := api.CreateUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    "invalid-email",
			Password: testutil.FakePassword(),
			Role:     "member",
		}
		var res api.CreateUserResponse
		httpRes, err := sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "email", "email")
	})

	t.Run("valid password", func(t *testing.T) {
		req := api.CreateUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: "no8char",
			Role:     "member",
		}
		var res api.CreateUserResponse
		httpRes, err := sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "password", "min")
	})

	t.Run("valid role", func(t *testing.T) {
		req := api.CreateUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: testutil.FakePassword(),
			Role:     "invalid-role",
		}

		var res api.CreateUserResponse
		httpRes, err := sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "role", "oneof")
	})

	t.Run("unique email", func(t *testing.T) {
		req := api.CreateUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: testutil.FakePassword(),
			Role:     "member",
		}

		var res api.CreateUserResponse
		httpRes, err := sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusCreated, httpRes.StatusCode)

		// Use a different username to avoid unique constraint violation. We only
		// care about email.
		req.Username = gofakeit.Username()
		httpRes, err = sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "email", "unique")
	})

	t.Run("unique username", func(t *testing.T) {
		req := api.CreateUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: testutil.FakePassword(),
			Role:     "member",
		}

		var res api.CreateUserResponse
		httpRes, err := sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusCreated, httpRes.StatusCode)

		// Use a different email to avoid unique constraint violation. We only care
		// about username.
		req.Email = gofakeit.Email()
		httpRes, err = sdk.CreateUser(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "username", "unique")
	})
}

func TestCreateUser_Success(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	req := api.CreateUserRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: testutil.FakePassword(),
		Role:     "member",
	}

	var res api.CreateUserResponse
	httpRes, err := sdk.CreateUser(req, &res)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	require.NotEmpty(t, res.Data.ID)
	require.Equal(t, req.Name, res.Data.Name)
	require.Equal(t, req.Username, res.Data.Username)
	require.Equal(t, req.Email, res.Data.Email)
	require.Equal(t, req.Role, res.Data.Role)
}

func TestCreateUser_OnlyAdminsCanCreateAdmins(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	memberReq := api.CreateUserRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: testutil.FakePassword(),
		Role:     "member",
	}
	var res api.CreateUserResponse
	httpRes, err := sdk.CreateUser(memberReq, &res)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	memberSDK := tEnv.AuthSDK(memberReq.Email, memberReq.Password)
	httpRes, err = memberSDK.CreateUser(api.CreateUserRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: testutil.FakePassword(),
		Role:     "admin",
	}, &res)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusForbidden, httpRes.StatusCode)
}

func TestDeleteUser_Success(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	req := api.CreateUserRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: testutil.FakePassword(),
		Role:     "member",
	}
	var res api.CreateUserResponse
	httpRes, err := sdk.CreateUser(req, &res)
	require.NoError(t, err, "error making create user request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	httpRes, err = sdk.DeleteUser(res.Data.ID)
	require.NoError(t, err, "error making delete user request")
	require.Equal(t, http.StatusNoContent, httpRes.StatusCode)
}

func TestDeleteUser_OnlyAdminsCanDelete(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	memberReq := api.CreateUserRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: testutil.FakePassword(),
		Role:     "member",
	}
	var res api.CreateUserResponse
	httpRes, err := sdk.CreateUser(memberReq, &res)
	require.NoError(t, err, "error making create user request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	memberSDK := tEnv.AuthSDK(memberReq.Email, memberReq.Password)
	httpRes, err = memberSDK.DeleteUser(res.Data.ID)
	require.NoError(t, err, "error making delete user request")
	require.Equal(t, http.StatusForbidden, httpRes.StatusCode)
}

func TestDeleteUser_CantSelfDelete(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	httpRes, err := sdk.DeleteUser(setup.Res().Data.ID)
	require.NoError(t, err, "error making delete user request")
	require.Equal(t, http.StatusForbidden, httpRes.StatusCode)
}
