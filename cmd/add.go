package main

import (
	"sync"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/struckchure/gv"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Download TypeScript definitions",
	Run: func(cmd *cobra.Command, args []string) {
		manager := gv.NewManager()
		var wg sync.WaitGroup
		for _, pkg := range args {
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				if err := manager.Install(p); err != nil {
					color.Red("❌ Error installing %s: %v\n", p, err)
				}
			}(pkg)
		}

		wg.Wait()
		color.Green("All .d.ts files downloaded!")
	},
	Args: cobra.MinimumNArgs(1),
}
