package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"memento/scripts"
	"memento/utils"
)

var typesCmd = &cobra.Command{
	Use:   "types",
	Short: "Manipulate or list the types configured in your memento project",
}

var typesAdd = &cobra.Command{
	Use:   "add",
	Short: "Add a type to memento project",

	Run: func(cmd *cobra.Command, args []string) {
		workingDir, err := utils.GetProjectPath()
		if err != nil { fmt.Println(err); return}

		scripts.Add(workingDir)
	},
}

var typesRemove = &cobra.Command{
	Use:   "remove",
	Short: "Remove a type from your memento project",

	Run: func(cmd *cobra.Command, args []string) {
		workingDir, err := utils.GetProjectPath()
		if err != nil { fmt.Println(err); return}

		scripts.Remove(workingDir)
	},
}

var typesList = &cobra.Command{
	Use:   "list",
	Short: "List all types in your memento project",

	Run: func(cmd *cobra.Command, args []string) {
		workingDir, err := utils.GetProjectPath()
		if err != nil { fmt.Println(err); return}

		scripts.List(workingDir)
	},
}

func init() {
	RootCmd.AddCommand(typesCmd)

	typesCmd.AddCommand(typesAdd)
	typesCmd.AddCommand(typesRemove)
	typesCmd.AddCommand(typesList)
}
