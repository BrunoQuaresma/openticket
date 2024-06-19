package api_test

import (
	"net/http"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestCreateComment_Success(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	t.Cleanup(tEnv.Close)
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	ticketReq := api.CreateTicketRequest{
		Title:       gofakeit.JobTitle(),
		Description: gofakeit.Sentence(10),
		Labels:      []string{gofakeit.HackerAbbreviation()},
	}
	var ticketRes api.CreateTicketResponse
	httpRes, err := sdk.CreateTicket(ticketReq, &ticketRes)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	commentReq := api.CreateCommentRequest{
		Content: gofakeit.Sentence(10),
	}
	var commentRes api.CreateCommentResponse
	httpRes, err = sdk.CreateComment(ticketRes.Data.ID, commentReq, &commentRes)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)
	require.NotEmpty(t, commentRes.Data.ID)
	require.Equal(t, commentReq.Content, commentRes.Data.Content)
}

func TestDeleteComment_Success(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	t.Cleanup(tEnv.Close)
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	ticketReq := api.CreateTicketRequest{
		Title:       gofakeit.JobTitle(),
		Description: gofakeit.Sentence(10),
		Labels:      []string{gofakeit.HackerAbbreviation()},
	}
	var ticketRes api.CreateTicketResponse
	httpRes, err := sdk.CreateTicket(ticketReq, &ticketRes)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	commentReq := api.CreateCommentRequest{
		Content: gofakeit.Sentence(10),
	}
	var commentRes api.CreateCommentResponse
	httpRes, err = sdk.CreateComment(ticketRes.Data.ID, commentReq, &commentRes)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	httpRes, err = sdk.DeleteComment(ticketRes.Data.ID, commentRes.Data.ID)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusNoContent, httpRes.StatusCode)
}

func TestDeleteComment_FailWhenUserIsNotAdminOrOwner(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	t.Cleanup(tEnv.Close)
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	ticketReq := api.CreateTicketRequest{
		Title:       gofakeit.JobTitle(),
		Description: gofakeit.Sentence(10),
		Labels:      []string{gofakeit.HackerAbbreviation()},
	}
	var ticketRes api.CreateTicketResponse
	httpRes, err := sdk.CreateTicket(ticketReq, &ticketRes)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	commentReq := api.CreateCommentRequest{
		Content: gofakeit.Sentence(10),
	}
	var commentRes api.CreateCommentResponse
	httpRes, err = sdk.CreateComment(ticketRes.Data.ID, commentReq, &commentRes)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	member := testutil.NewMember(t, sdk)
	memberSdk := tEnv.AuthSDK(member.Email, member.Password)

	httpRes, err = memberSdk.DeleteComment(ticketRes.Data.ID, commentRes.Data.ID)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusForbidden, httpRes.StatusCode)
}

func TestPatchComment_Success(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	t.Cleanup(tEnv.Close)
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	ticketReq := api.CreateTicketRequest{
		Title:       gofakeit.JobTitle(),
		Description: gofakeit.Sentence(10),
		Labels:      []string{gofakeit.HackerAbbreviation()},
	}
	var ticketRes api.CreateTicketResponse
	httpRes, err := sdk.CreateTicket(ticketReq, &ticketRes)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	commentReq := api.CreateCommentRequest{
		Content: gofakeit.Sentence(10),
	}
	var commentRes api.CreateCommentResponse
	httpRes, err = sdk.CreateComment(ticketRes.Data.ID, commentReq, &commentRes)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	newContent := gofakeit.Sentence(10)
	patchReq := api.PatchCommentRequest{
		Content: newContent,
	}
	var patchRes api.PatchCommentResponse
	httpRes, err = sdk.PatchComment(ticketRes.Data.ID, commentRes.Data.ID, patchReq, &patchRes)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.Equal(t, newContent, patchRes.Data.Content)
}

func TestPatchComment_FailWhenUserIsNotAdminOrOwner(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	t.Cleanup(tEnv.Close)
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	ticketReq := api.CreateTicketRequest{
		Title:       gofakeit.JobTitle(),
		Description: gofakeit.Sentence(10),
		Labels:      []string{gofakeit.HackerAbbreviation()},
	}
	var ticketRes api.CreateTicketResponse
	httpRes, err := sdk.CreateTicket(ticketReq, &ticketRes)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	commentReq := api.CreateCommentRequest{
		Content: gofakeit.Sentence(10),
	}
	var commentRes api.CreateCommentResponse
	httpRes, err = sdk.CreateComment(ticketRes.Data.ID, commentReq, &commentRes)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusCreated, httpRes.StatusCode)

	member := testutil.NewMember(t, sdk)
	memberSdk := tEnv.AuthSDK(member.Email, member.Password)

	newContent := gofakeit.Sentence(10)
	patchReq := api.PatchCommentRequest{
		Content: newContent,
	}
	var patchRes api.PatchCommentResponse
	httpRes, err = memberSdk.PatchComment(ticketRes.Data.ID, commentRes.Data.ID, patchReq, &patchRes)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusForbidden, httpRes.StatusCode)
}
