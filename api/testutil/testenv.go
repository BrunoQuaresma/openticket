package testutil

import (
	"io"
	"net"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/BrunoQuaresma/openticket/sdk"
	"github.com/brianvoe/gofakeit"
)

type TestEnv struct {
	database *TestDatabase
	server   *api.Server
	t        *testing.T
}

func NewEnv(t *testing.T) *TestEnv {
	tEnv := TestEnv{t: t}

	dbPort, err := getFreePort()
	if err != nil {
		t.Fatal("error getting free port for db: " + err.Error())
	}
	tEnv.database = NewTestDatabase(dbPort)

	port, err := getFreePort()
	if err != nil {
		t.Fatal("error getting free port for server: " + err.Error())
	}
	tEnv.server = api.NewServer(api.ServerOptions{
		DatabaseURL: tEnv.database.URL(),
		Mode:        api.TestMode,
		Port:        port,
	})
	return &tEnv
}

func (tEnv *TestEnv) Start() {
	err := tEnv.database.Start(io.Discard)
	if err != nil {
		tEnv.t.Fatal("error starting test database: " + err.Error())
	}

	tEnv.server.Start()
}

func (tEnv *TestEnv) Close() {
	tEnv.database.Stop()
	tEnv.server.Close()
}

func (tEnv *TestEnv) Server() *api.Server {
	return tEnv.server
}

func (tEnv *TestEnv) Setup() api.SetupRequest {
	req := api.SetupRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: FakePassword(),
	}
	var res api.Response[any]
	sdk := tEnv.SDK()
	_, err := sdk.Setup(req, &res)
	if err != nil {
		tEnv.t.Fatal("error making setup request: " + err.Error())
	}
	return req
}

func (tEnv *TestEnv) DBQueries() *database.Queries {
	return tEnv.server.DBQueries()
}

func (tEnv *TestEnv) SDK() *sdk.Client {
	return sdk.New(tEnv.Server().URL())
}

func (tEnv *TestEnv) AuthSDK(email string, password string) *sdk.Client {
	sdk := tEnv.SDK()
	var loginRes api.LoginResponse
	_, err := sdk.Login(api.LoginRequest(api.LoginRequest{
		Email:    email,
		Password: password,
	}), &loginRes)
	if err != nil {
		tEnv.t.Fatal("error making login request" + err.Error())
	}
	sdk.Authenticate(loginRes.Data.SessionToken)
	return sdk
}

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}
