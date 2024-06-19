package api_test

import (
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestAPI_CreateUser(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	t.Run("error: no authentication", func(t *testing.T) {
		t.Parallel()

		noAuthSDK := tEnv.SDK()
		httpRes, err := noAuthSDK.CreateUser(api.CreateUserRequest{}, &api.CreateUserResponse{})
		require.NoError(t, err, "error making create user request")
		require.Equal(t, http.StatusUnauthorized, httpRes.StatusCode)
	})

	t.Run("validation", func(t *testing.T) {
		t.Parallel()

		t.Run("error: missing required fields", func(t *testing.T) {
			t.Parallel()

			var (
				req api.CreateUserRequest
				res api.CreateUserResponse
			)
			httpRes, err := sdk.CreateUser(req, &res)
			require.NoError(t, err, "error making request")
			require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
			testutil.RequireValidationError(t, res.Errors, "name", "required")
			testutil.RequireValidationError(t, res.Errors, "username", "required")
			testutil.RequireValidationError(t, res.Errors, "email", "required")
			testutil.RequireValidationError(t, res.Errors, "password", "required")
			testutil.RequireValidationError(t, res.Errors, "role", "required")
		})

		t.Run("error: invalid email", func(t *testing.T) {
			t.Parallel()

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

		t.Run("error: invalid password", func(t *testing.T) {
			t.Parallel()

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

		t.Run("error: invalid role", func(t *testing.T) {
			t.Parallel()

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

		t.Run("error: duplicated email", func(t *testing.T) {
			t.Parallel()

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

		t.Run("error: duplicated username", func(t *testing.T) {
			t.Parallel()

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
	})

	t.Run("success", func(t *testing.T) {
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
	})

	t.Run("error: non admins creating users", func(t *testing.T) {
		t.Parallel()

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
	})
}

func TestAPI_DeleteUser(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

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
	})

	t.Run("error: non admins deleting users", func(t *testing.T) {
		t.Parallel()

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
	})

	t.Run("error: deleting self", func(t *testing.T) {
		t.Parallel()

		httpRes, err := sdk.DeleteUser(setup.Res().Data.ID)
		require.NoError(t, err, "error making delete user request")
		require.Equal(t, http.StatusForbidden, httpRes.StatusCode)
	})
}

func TestAPI_PatchUser(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	t.Run("patching single fields", func(t *testing.T) {
		t.Parallel()

		t.Run("success: name", func(t *testing.T) {
			t.Parallel()

			createMemberReq := api.CreateUserRequest{
				Name:     gofakeit.Name(),
				Username: gofakeit.Username(),
				Email:    gofakeit.Email(),
				Role:     "member",
				Password: testutil.FakePassword(),
			}
			var memberUserRes api.CreateUserResponse
			httpRes, err := sdk.CreateUser(createMemberReq, &memberUserRes)
			require.NoError(t, err, "error making create user request")
			require.Equal(t, http.StatusCreated, httpRes.StatusCode)

			patchUserReq := api.PatchUserRequest{
				Name: gofakeit.Name(),
			}
			var updatedUserRes api.PatchUserResponse
			httpRes, err = sdk.PatchUser(memberUserRes.Data.ID, patchUserReq, &updatedUserRes)
			require.NoError(t, err, "error making patch user request")
			require.Equal(t, http.StatusOK, httpRes.StatusCode)

			require.Equal(t, patchUserReq.Name, updatedUserRes.Data.Name)
			require.Equal(t, memberUserRes.Data.Username, updatedUserRes.Data.Username)
			require.Equal(t, memberUserRes.Data.Email, updatedUserRes.Data.Email)
			require.Equal(t, memberUserRes.Data.Role, updatedUserRes.Data.Role)
		})

		t.Run("success: username", func(t *testing.T) {
			t.Parallel()

			createMemberReq := api.CreateUserRequest{
				Name:     gofakeit.Name(),
				Username: gofakeit.Username(),
				Email:    gofakeit.Email(),
				Role:     "member",
				Password: testutil.FakePassword(),
			}
			var memberUserRes api.CreateUserResponse
			httpRes, err := sdk.CreateUser(createMemberReq, &memberUserRes)
			require.NoError(t, err, "error making create user request")
			require.Equal(t, http.StatusCreated, httpRes.StatusCode)

			patchUserReq := api.PatchUserRequest{
				Username: gofakeit.Username(),
			}
			var updatedUserRes api.PatchUserResponse
			httpRes, err = sdk.PatchUser(memberUserRes.Data.ID, patchUserReq, &updatedUserRes)
			require.NoError(t, err, "error making patch user request")
			require.Equal(t, http.StatusOK, httpRes.StatusCode)

			require.Equal(t, memberUserRes.Data.Name, updatedUserRes.Data.Name)
			require.Equal(t, patchUserReq.Username, updatedUserRes.Data.Username)
			require.Equal(t, memberUserRes.Data.Email, updatedUserRes.Data.Email)
			require.Equal(t, memberUserRes.Data.Role, updatedUserRes.Data.Role)
		})

		t.Run("success: email", func(t *testing.T) {
			t.Parallel()

			createMemberReq := api.CreateUserRequest{
				Name:     gofakeit.Name(),
				Username: gofakeit.Username(),
				Email:    gofakeit.Email(),
				Role:     "member",
				Password: testutil.FakePassword(),
			}
			var memberUserRes api.CreateUserResponse
			httpRes, err := sdk.CreateUser(createMemberReq, &memberUserRes)
			require.NoError(t, err, "error making create user request")
			require.Equal(t, http.StatusCreated, httpRes.StatusCode)

			patchUserReq := api.PatchUserRequest{
				Email: gofakeit.Email(),
			}
			var updatedUserRes api.PatchUserResponse
			httpRes, err = sdk.PatchUser(memberUserRes.Data.ID, patchUserReq, &updatedUserRes)
			require.NoError(t, err, "error making patch user request")
			require.Equal(t, http.StatusOK, httpRes.StatusCode)

			require.Equal(t, memberUserRes.Data.Name, updatedUserRes.Data.Name)
			require.Equal(t, memberUserRes.Data.Username, updatedUserRes.Data.Username)
			require.Equal(t, patchUserReq.Email, updatedUserRes.Data.Email)
			require.Equal(t, memberUserRes.Data.Role, updatedUserRes.Data.Role)
		})
	})

	t.Run("success: admins patching other users", func(t *testing.T) {
		t.Parallel()

		createUserReq := api.CreateUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Role:     "member",
			Password: testutil.FakePassword(),
		}
		var newUserRes api.CreateUserResponse
		httpRes, err := sdk.CreateUser(createUserReq, &newUserRes)
		require.NoError(t, err, "error making create user request")
		require.Equal(t, http.StatusCreated, httpRes.StatusCode)

		patchUserReq := api.PatchUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Role:     "admin",
		}
		var updatedUserRes api.PatchUserResponse
		httpRes, err = sdk.PatchUser(newUserRes.Data.ID, patchUserReq, &updatedUserRes)
		require.NoError(t, err, "error making patch user request")
		require.Equal(t, http.StatusOK, httpRes.StatusCode)

		require.Equal(t, patchUserReq.Name, updatedUserRes.Data.Name)
		require.Equal(t, patchUserReq.Username, updatedUserRes.Data.Username)
		require.Equal(t, patchUserReq.Email, updatedUserRes.Data.Email)
		require.Equal(t, patchUserReq.Role, updatedUserRes.Data.Role)
	})

	t.Run("success: members patching their own information", func(t *testing.T) {
		t.Parallel()

		createMemberReq := api.CreateUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Role:     "member",
			Password: testutil.FakePassword(),
		}
		var memberUserRes api.CreateUserResponse
		httpRes, err := sdk.CreateUser(createMemberReq, &memberUserRes)
		require.NoError(t, err, "error making create user request")
		require.Equal(t, http.StatusCreated, httpRes.StatusCode)

		patchMemberReq := api.PatchUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
		}
		var updatedUserRes api.PatchUserResponse
		memberSDK := tEnv.AuthSDK(createMemberReq.Email, createMemberReq.Password)
		httpRes, err = memberSDK.PatchUser(memberUserRes.Data.ID, patchMemberReq, &updatedUserRes)
		require.NoError(t, err, "error making patch user request")
		require.Equal(t, http.StatusOK, httpRes.StatusCode)

		require.Equal(t, patchMemberReq.Name, updatedUserRes.Data.Name)
		require.Equal(t, patchMemberReq.Username, updatedUserRes.Data.Username)
		require.Equal(t, patchMemberReq.Email, updatedUserRes.Data.Email)
		require.Equal(t, "member", updatedUserRes.Data.Role)

		httpRes, err = memberSDK.PatchUser(setup.Res().Data.ID, api.PatchUserRequest{}, &api.PatchUserResponse{})
		require.NoError(t, err, "error making patch user request")
		require.Equal(t, http.StatusForbidden, httpRes.StatusCode)
	})

	t.Run("error: non admins patching roles", func(t *testing.T) {
		t.Parallel()

		createMemberReq := api.CreateUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Role:     "member",
			Password: testutil.FakePassword(),
		}
		var memberUserRes api.CreateUserResponse
		httpRes, err := sdk.CreateUser(createMemberReq, &memberUserRes)
		require.NoError(t, err, "error making create user request")
		require.Equal(t, http.StatusCreated, httpRes.StatusCode)

		patchMemberReq := api.PatchUserRequest{
			Role: "admin",
		}
		var updatedUserRes api.PatchUserResponse
		memberSDK := tEnv.AuthSDK(createMemberReq.Email, createMemberReq.Password)
		httpRes, err = memberSDK.PatchUser(memberUserRes.Data.ID, patchMemberReq, &updatedUserRes)
		require.NoError(t, err, "error making patch user request")
		require.Equal(t, http.StatusForbidden, httpRes.StatusCode)
	})

	t.Run("success: admins patching roles", func(t *testing.T) {
		t.Parallel()

		createMemberReq := api.CreateUserRequest{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Role:     "member",
			Password: testutil.FakePassword(),
		}
		var memberUserRes api.CreateUserResponse
		httpRes, err := sdk.CreateUser(createMemberReq, &memberUserRes)
		require.NoError(t, err, "error making create user request")
		require.Equal(t, http.StatusCreated, httpRes.StatusCode)

		patchMemberReq := api.PatchUserRequest{
			Role: "admin",
		}
		var updatedUserRes api.PatchUserResponse
		httpRes, err = sdk.PatchUser(memberUserRes.Data.ID, patchMemberReq, &updatedUserRes)
		require.NoError(t, err, "error making patch user request")
		require.Equal(t, http.StatusOK, httpRes.StatusCode)
		require.Equal(t, "admin", updatedUserRes.Data.Role)
	})

	t.Run("error: demoting single admin", func(t *testing.T) {
		t.Parallel()

		httpRes, err := sdk.PatchUser(
			setup.Res().Data.ID,
			api.PatchUserRequest{Role: "member"},
			&api.PatchUserResponse{},
		)
		require.NoError(t, err, "error making patch user request")
		require.Equal(t, http.StatusForbidden, httpRes.StatusCode)
	})
}
