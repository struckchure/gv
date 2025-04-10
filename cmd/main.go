package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	// Root command
	var rootCmd = &cobra.Command{Use: "gv"}

	// Add the 'add' command to root command
	rootCmd.AddCommand(addCommand)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
