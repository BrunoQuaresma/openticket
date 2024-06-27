package api_test

import (
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestAPI_CreateAssignment(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	var userRes api.CreateUserResponse
	httpRes, err := sdk.CreateUser(api.CreateUserRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: testutil.FakePassword(),
		Role:     "member",
	}, &userRes)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	var ticketRes api.CreateTicketResponse
	_, err = sdk.CreateTicket(api.CreateTicketRequest{
		Title:       gofakeit.Job().Title,
		Description: gofakeit.Sentence(10),
	}, &ticketRes)
	require.NoError(t, err, "error creating ticket")

	t.Run("success: create assignment", func(t *testing.T) {
		t.Parallel()

		var assignmentRes api.CreateAssignmentResponse
		httpRes, err = sdk.CreateAssignment(ticketRes.Data.ID, userRes.Data.ID, &assignmentRes)
		require.NoError(t, err, "error creating assignment")
		require.Equal(t, http.StatusCreated, httpRes.StatusCode)
		require.Equal(t, userRes.Data.ID, assignmentRes.Data.UserID)
		require.Equal(t, ticketRes.Data.ID, assignmentRes.Data.TicketID)
		require.NotEmpty(t, assignmentRes.Data.ID)
	})

	t.Run("success: delete assignment", func(t *testing.T) {
		t.Parallel()

		_, err = sdk.DeleteAssignment(1)
		require.NoError(t, err, "error deleting assignment")
	})
}
