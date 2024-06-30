package main

import (
	"fmt"
	"os"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/spf13/cobra"
)

func main() {

	rootCmd := &cobra.Command{
		Use:   "openticket",
		Short: "Openticket is a ticketing system for managing tickets.",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetInt("port")
			server := api.NewServer(api.ServerOptions{
				Mode: api.ProductionMode,
				Port: port,
			})
			defer server.Close()
			server.Start()
		},
	}

	rootCmd.Flags().IntP("port", "p", 8080, "Port to run the server on.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
