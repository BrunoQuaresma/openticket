package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/api/database"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "openticket",
		Short: "Openticket is a ticketing system for managing tickets.",
		Run: func(cmd *cobra.Command, args []string) {
			s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)

			loading(s, "Starting local database...", "✔ Local database started")
			localDB := database.NewLocalDatabase(5432, ".openticket", io.Discard)
			defer localDB.Stop()
			err := localDB.Start()
			if err != nil {
				log.Fatal("error starting local database: " + err.Error())
			}
			s.Stop()

			loading(s, "Applying migrations...", "✔ Migrations applied")
			err = localDB.Migrate()
			if err != nil {
				log.Fatal("error migrating local database: " + err.Error())
			}
			s.Stop()

			loading(s, "Connecting to the database...", "✔ Database connected")
			conn, err := database.Connect(localDB.URL())
			if err != nil {
				log.Fatal("error connecting to database: " + err.Error())
			}
			defer conn.Close()
			s.Stop()

			port, _ := cmd.Flags().GetInt("port")
			server := api.NewServer(port, &conn, api.ProductionMode)
			loading(s, "Starting server...", "✔ Server started on "+server.URL())
			go func() {
				defer server.Close()
				server.Start()
			}()
			s.Stop()

			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-c
				server.Close()
				conn.Close()
				localDB.Stop()
				os.Exit(1)
			}()

			select {}
		},
	}

	rootCmd.Flags().IntP("port", "p", 8080, "Port to run the server on.")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("error executing command: " + err.Error())
	}
}

func loading(s *spinner.Spinner, message string, success string) {
	s.FinalMSG = success + "\n"
	s.Suffix = " " + message
	s.Start()
}
