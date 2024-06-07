package api_test

import (
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/stretchr/testify/require"
)

func TestCreateTicket_Success(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	req := api.CreateTicketRequest{
		Title:       "User cannot login",
		Description: "User cannot login to the system",
		Labels:      []string{"bug", "customer", "login"},
	}

	var res api.CreateTicketResponse
	httpRes, err := sdk.CreateTicket(req, &res)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	require.NotEmpty(t, res.Data.ID)
	require.Equal(t, req.Title, res.Data.Title)
	require.Equal(t, req.Description, res.Data.Description)
	require.Equal(t, req.Labels, res.Data.Labels)
	require.Equal(t, setup.Res().Data.ID, res.Data.CreatedBy.ID)
}

func TestCreateTicket_Validation(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	t.Run("required fields", func(t *testing.T) {
		req := api.CreateTicketRequest{}
		var res api.CreateTicketResponse
		httpRes, err := sdk.CreateTicket(req, &res)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusBadRequest, httpRes.StatusCode)
		testutil.RequireValidationError(t, res.Errors, "title", "required")
		testutil.RequireValidationError(t, res.Errors, "description", "required")
	})
}
