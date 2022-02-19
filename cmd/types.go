package cmd

import (
	"github.com/spf13/cobra"
	"memento/scripts"
)

var typesCmd = &cobra.Command{
	Use:   "types",
	Short: "Manipulate or list the types configured in your memento project",
}

var typesAdd = &cobra.Command{
	Use:   "add",
	Short: "Add a type to memento project",

	Run: func(cmd *cobra.Command, args []string) {
		scripts.Add()
	},
}

var typesRemove = &cobra.Command{
	Use:   "remove",
	Short: "Remove a type from your memento project",

	Run: func(cmd *cobra.Command, args []string) {
		scripts.Remove()
	},
}

var typesList = &cobra.Command{
	Use:   "list",
	Short: "List all types in your memento project",

	Run: func(cmd *cobra.Command, args []string) {
		scripts.List()
	},
}

func init() {
	RootCmd.AddCommand(typesCmd)

	typesCmd.AddCommand(typesAdd)
	typesCmd.AddCommand(typesRemove)
	typesCmd.AddCommand(typesList)
}
