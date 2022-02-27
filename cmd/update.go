package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"memento/scripts"
	"memento/utils"
	"os"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all permanent datasources in your project",
	Long: `Update all the permanent datasources in your memento project. This is automatically done on every startup
by a background service, but this command can be used to do it manually.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Too many arguments provided. Check memento update --help")
		} else {
			workingDir, err := os.Getwd()
			if utils.Handle(err) != nil {
				return
			}

			scripts.Update(workingDir)
		}

	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
