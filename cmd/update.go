package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"memento/scripts"
	"memento/utils"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all permanent datasources in your project",
	Long: `Update all the permanent datasources in your memento project. This is automatically done on every startup
by a background service, but this command can be used to do it manually.`,
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")

		if all {
			scripts.Update("all-console")
		} else {
			workingDir, err := utils.GetProjectPath()
			if err != nil { fmt.Println(err); return}

			if utils.Handle(err) != nil {
				return
			}

			scripts.Update(workingDir)
		}

	},
}

func init() {
	RootCmd.AddCommand(updateCmd)

	updateCmd.Flags().Bool("all", false, "Update all permanent data sources in all your project")
}
