package testutils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
)

type TestDatabase struct {
	Username string
	Password string
	Database string
	Port     uint32
	Conn     *embeddedpostgres.EmbeddedPostgres
}

func (testDB *TestDatabase) Start(logger io.Writer) error {
	testDB.Conn = embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Port(testDB.Port).
			Username(testDB.Username).
			Password(testDB.Password).
			Database(testDB.Database).
			Logger(logger),
	)
	err := testDB.Conn.Start()
	if err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	go func() {
		<-stop
		err := testDB.Stop()
		if err != nil {
			panic("error shutting down database: " + err.Error())
		} else {
			log.Println("database gracefully stopped")
		}
	}()

	migrateCmd := exec.Command("./scripts/migrate.sh")
	migrateCmd.Env = append(migrateCmd.Env, "POSTGRES_DB_URL="+testDB.URL())
	migrateCmd.Stdout = logger
	migrateCmd.Stderr = logger
	err = migrateCmd.Run()
	if err != nil {
		testDB.Stop()
		return errors.New("error running migrations: " + err.Error())
	}

	return nil
}

func (testDB *TestDatabase) Stop() error {
	return testDB.Conn.Stop()
}

func (testDB *TestDatabase) URL() string {
	return "postgresql://" + testDB.Username + ":" + testDB.Password + "@localhost:" + fmt.Sprint(testDB.Port) + "/" + testDB.Database + "?sslmode=disable"
}

func NewTestDatabase() *TestDatabase {
	port, err := getFreePort()
	if err != nil {
		panic("error getting free port: " + err.Error())
	}

	return &TestDatabase{
		Username: "postgres",
		Password: "postgres",
		Database: "postgres",
		Port:     uint32(port),
	}
}
