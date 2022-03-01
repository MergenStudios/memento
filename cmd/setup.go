package cmd

import (
	"github.com/spf13/cobra"
	"memento/scripts"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the directory structure needed by memento",
	Long: `The setup command sets up the directory structure needed by memento,
including a predefined configuration of types for the user to use.`,

	Run: func(cmd *cobra.Command, args []string) {
		scripts.Setup()
	},
}

func init() {
	RootCmd.AddCommand(setupCmd)
}
