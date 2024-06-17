package api_test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/testutil"
	"github.com/brianvoe/gofakeit/v7"
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

func TestTickets_Success(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	numberOfTickets := 5
	for i := range numberOfTickets {
		var res api.CreateTicketResponse
		httpRes, err := sdk.CreateTicket(api.CreateTicketRequest{
			Title:       gofakeit.JobTitle(),
			Description: gofakeit.HackerPhrase(),
			Labels:      []string{gofakeit.HackerAbbreviation()},
		}, &res)
		require.NoError(t, err, "error on create ticket request")
		require.Equal(t, http.StatusCreated, httpRes.StatusCode, "error creating ticket "+fmt.Sprint(i))
	}

	var res api.TicketsResponse
	httpRes, err := sdk.Tickets(&res, nil)
	require.NoError(t, err, "error making request")
	require.Equal(t, http.StatusOK, httpRes.StatusCode)
	require.Len(t, res.Data, numberOfTickets)
}

func TestTickets_FilterByLabel(t *testing.T) {
	t.Parallel()

	tEnv := testutil.NewEnv(t)
	tEnv.Start()
	defer tEnv.Close()
	setup := tEnv.Setup()
	sdk := tEnv.AuthSDK(setup.Req().Email, setup.Req().Password)

	var i int
	createTicket := func(title string, labels []string) {
		var res api.CreateTicketResponse
		httpRes, err := sdk.CreateTicket(api.CreateTicketRequest{
			Title:       title,
			Description: gofakeit.HackerPhrase(),
			Labels:      labels,
		}, &res)
		require.NoError(t, err, "error on create ticket request")
		require.Equal(t, http.StatusCreated, httpRes.StatusCode, "error creating ticket with labels "+strings.Join(labels, ",")+" and index "+fmt.Sprint(i))
		i++
	}

	createTicket("login not working", []string{"bug"})
	createTicket("register not working", []string{"bug"})
	createTicket("wrong validation", []string{"bug"})
	createTicket("add search", []string{"feature", "site"})
	createTicket("add validation", []string{"feature", "site"})
	createTicket("add health endpoint", []string{"feature", "api"})
	createTicket("be able to config using terraform", []string{"request"})

	t.Run("one label", func(t *testing.T) {
		urlValues := url.Values{
			"q": []string{"label:bug"},
		}
		var res api.TicketsResponse
		httpRes, err := sdk.Tickets(&res, &urlValues)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusOK, httpRes.StatusCode, "error getting tickets")
		require.Len(t, res.Data, 3)
	})

	t.Run("multiple labels as OR", func(t *testing.T) {
		urlValues := url.Values{
			"q": []string{"label:bug,request"},
		}
		var res api.TicketsResponse
		httpRes, err := sdk.Tickets(&res, &urlValues)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusOK, httpRes.StatusCode, "error getting tickets")
		require.Len(t, res.Data, 4)
	})

	t.Run("multiple labels as AND", func(t *testing.T) {
		urlValues := url.Values{
			"q": []string{"label:feature label:site"},
		}
		var res api.TicketsResponse
		httpRes, err := sdk.Tickets(&res, &urlValues)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusOK, httpRes.StatusCode, "error getting tickets")
		require.Len(t, res.Data, 2)
	})

	t.Run("search by title", func(t *testing.T) {
		urlValues := url.Values{
			"q": []string{"not working"},
		}
		var res api.TicketsResponse
		httpRes, err := sdk.Tickets(&res, &urlValues)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusOK, httpRes.StatusCode, "error getting tickets")
		require.Len(t, res.Data, 2)
	})

	t.Run("search by title and label", func(t *testing.T) {
		urlValues := url.Values{
			"q": []string{"validation label:site"},
		}
		var res api.TicketsResponse
		httpRes, err := sdk.Tickets(&res, &urlValues)
		require.NoError(t, err, "error making request")
		require.Equal(t, http.StatusOK, httpRes.StatusCode, "error getting tickets")
		require.Len(t, res.Data, 1)
	})
}
