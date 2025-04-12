package main

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/struckchure/gv"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Download TypeScript definitions",
	Run:   func(cmd *cobra.Command, args []string) { addService(args) },
	Args:  cobra.MinimumNArgs(1),
}

func addService(pkgs []string) {
	manager := gv.NewManager()
	manager.Install(pkgs...)
	color.Green("All .d.ts files downloaded!")
}

func init() {
	rootCmd.AddCommand(addCommand)
}
