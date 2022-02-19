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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
