package testutil

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/sdk"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func NewMember(t *testing.T, sdk *sdk.Client) (api.CreateUserRequest, api.CreateUserResponse) {
	req := api.CreateUserRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: FakePassword(),
		Role:     "member",
	}

	var res api.CreateUserResponse
	httpRes, err := sdk.CreateUser(req, &res)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	return req, res
}

func RequireValidationError(t *testing.T, errors []api.ValidationError, field string, validator string) {
	hasError := false

	if len(errors) == 0 {
		hasError = false
	}

	for _, e := range errors {
		if e.Field == field && e.Validator == validator {
			hasError = true
			break
		}
	}

	require.True(t, hasError, fmt.Sprintf("expected error for field %s with validator %s", field, validator))
}

func FakePassword() string {
	return gofakeit.Password(true, true, true, true, false, 15)
}
