package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-get-cli",
	Short: "CLI Tools for go get to manage and find pacakges",
	Long:  "CLI Tools for go get to manage and find pacakges",
}

func init() {
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
