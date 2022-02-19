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
	Short: "Import a single or multiple datapoints to your project.",

	Run: func(cmd *cobra.Command, args []string) {
		// check if the args are valid
		if len(args) != 2 {
			if len(args) == 0 {
				fmt.Println("No arguments provided. Check memento import --help for more information.")
				return
			} else if len(args) == 2 {
				fmt.Println("Not enough arguments provided. Check memento import --help for more information.")
				return
			} else if len(args) > 2 {
				fmt.Println("Too many arguments provided. Check memento import --help for more information.")
			}
		} else if !utils.IsType(args[0]) {
			fmt.Printf("Uknown datatype: %s. Check memento types list for all types.", args[0])
		} else if _, err := os.Stat(args[1]); os.IsNotExist(err) {
			fmt.Printf("No such file or directory: %s", args[1])
		} else {
			datatype := args[0]
			path := args[1]

			err := scripts.ImportDatapoints(datatype, path)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(importCmd)
}
