package main

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/struckchure/gv"
)

var syncCommand = &cobra.Command{
	Use:   "sync",
	Short: "Sync Dependencies",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			color.Red(err.Error())
			return
		}

		manager := gv.NewManager(gv.ManagerOptions{ConfigFile: configFile})
		manager.Sync(configFile)
	},
}

func init() {
	syncCommand.Flags().StringP("config", "c", "./config.yaml", "Config file")

	rootCmd.AddCommand(syncCommand)
}
