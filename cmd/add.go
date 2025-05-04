package main

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/struckchure/gv"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Update config file package registry and install typescript definitions.",
	Run:   addService,
	Args:  cobra.MinimumNArgs(1),
}

func addService(cmd *cobra.Command, packages []string) {
	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		color.Red(err.Error())
		return
	}

	manager := gv.NewManager(gv.ManagerOptions{ConfigFile: configFile})
	manager.Add(packages...)
}

func init() {
	addCommand.Flags().StringP("config", "c", "./config.yaml", "Config file")

	rootCmd.AddCommand(addCommand)
}
