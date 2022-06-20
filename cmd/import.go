package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"memento/scripts"
	"memento/utils"
	"os"
)

var importCmd = &cobra.Command{
	Use:   "import [flags] [PATH]",
	Short: "Import datapoints to your memento project",
	Long:  `This command is used to import datapoints into your project.`,

	Run: func(cmd *cobra.Command, args []string) {
		workingDir, err := utils.GetProjectPath()
		if err != nil {
			fmt.Println(err)
			return
		}

		// check if the args are valid
		if len(args) != 1 {
			if len(args) == 0 {
				fmt.Println("No arguments provided. Check memento import --help for more information")
				return
			} else if len(args) > 1 {
				fmt.Println("Too many arguments provided. Check memento import --help for more information")
			}
		} else if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			fmt.Printf("No such file or directory: %s", args[1])
		} else {
			path := args[0]

			scripts.ImportDatapoints(path, workingDir)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(importCmd)
}
