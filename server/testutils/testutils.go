package testutils

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	server "github.com/BrunoQuaresma/openticket"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
)

type TestDatabaseConfig struct {
	Username string
	Password string
	Database string
	Port     uint32
}

func NewTestDatabaseConfig() TestDatabaseConfig {
	return TestDatabaseConfig{
		Username: "postgres",
		Password: "postgres",
		Database: "postgres",
		Port:     5433,
	}
}

type TestServerConfig struct {
	Database TestDatabaseConfig
	Debug    bool
}

type Stop func()

func RunTestServer(c TestServerConfig) Stop {
	logger := io.Discard
	if c.Debug {
		logger = os.Stdout
	}

	pg := embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Port(c.Database.Port).
			Username(c.Database.Username).
			Password(c.Database.Password).
			Database(c.Database.Database).
			Logger(logger),
	)
	err := pg.Start()
	if err != nil {
		panic("error starting postgres: " + err.Error())
	}
	dbURL := "postgresql://" + c.Database.Username + ":" + c.Database.Password + "@localhost:" + fmt.Sprint(c.Database.Port) + "/" + c.Database.Database + "?sslmode=disable"

	migrateCmd := exec.Command("./scripts/migrate.sh")
	migrateCmd.Env = append(migrateCmd.Env, "POSTGRES_DB_URL="+dbURL)
	migrateCmd.Stdout = logger
	migrateCmd.Stderr = logger
	err = migrateCmd.Run()
	if err != nil {
		pg.Stop()
		panic("error running migrations: " + err.Error())
	}

	s := server.Start(server.Options{
		DatabaseURL: dbURL,
		Debug:       c.Debug,
	})

	return func() {
		pg.Stop()
		s.Close()
	}
}
