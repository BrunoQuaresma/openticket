package testutils

import (
	"fmt"
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

type Stop func()

func RunTestServer(dbConfig TestDatabaseConfig) Stop {
	pg := embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Port(dbConfig.Port).
			Username(dbConfig.Username).
			Password(dbConfig.Password).
			Database(dbConfig.Database),
	)
	err := pg.Start()
	if err != nil {
		panic("error starting postgres: " + err.Error())
	}
	dbURL := "postgresql://" + dbConfig.Username + ":" + dbConfig.Password + "@localhost:" + fmt.Sprint(dbConfig.Port) + "/" + dbConfig.Database + "?sslmode=disable"

	migrateCmd := exec.Command("./scripts/migrate.sh")
	migrateCmd.Env = append(migrateCmd.Env, "POSTGRES_DB_URL="+dbURL)
	migrateCmd.Stdout = os.Stdout
	migrateCmd.Stderr = os.Stderr
	err = migrateCmd.Run()
	if err != nil {
		pg.Stop()
		panic("error running migrations: " + err.Error())
	}

	s := server.Start(server.Options{
		DatabaseURL: dbURL,
	})

	return func() {
		pg.Stop()
		s.Close()
	}
}
