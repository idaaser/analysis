package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	execute()
}

var rootCmd = &cobra.Command{
	Use:   "idaas analysis tool help",
	Short: "how meeting uses idaas",
	Long:  "how meeting uses idaas",
}

func execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
