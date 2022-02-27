package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"memento/scripts"
	"memento/utils"
	"os"
)

var importCmd = &cobra.Command{
	Use:   "import [flags] [TYPE] [PATH]",
	Short: "Import datapoints to your memento project",
	Long: `This command is used to import datapoints into your project. It takes in a type and a path as arguments.
You can manipulate the types in your memento project with the types command, 
check memento types --help for more information.`,

	Run: func(cmd *cobra.Command, args []string) {
		workingDir, _ := os.Getwd()

		// check if the args are valid
		if len(args) != 2 {
			if len(args) == 0 {
				fmt.Println("No arguments provided. Check memento import --help for more information")
				return
			} else if len(args) == 2 {
				fmt.Println("Not enough arguments provided. Check memento import --help for more information")
				return
			} else if len(args) > 2 {
				fmt.Println("Too many arguments provided. Check memento import --help for more information")
			}
		} else if ok, _ := utils.IsType(args[0], workingDir); !ok {
			fmt.Println("Args: ", workingDir, args[0])
			fmt.Printf("Uknown datatype: %s. Check memento types list for all types.", args[0])
		} else if _, err := os.Stat(args[1]); os.IsNotExist(err) {
			fmt.Printf("No such file or directory: %s", args[1])
		} else {
			datatype := args[0]
			path := args[1]
			permanent, _ := cmd.Flags().GetBool("permanent")

			fmt.Println(workingDir)

			scripts.ImportDatapoints(datatype, path, workingDir, permanent, true)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(importCmd)

	importCmd.Flags().Bool("permanent", false, "Add this directory to be checked for new files on every startup.")
}
