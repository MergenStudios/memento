package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"memento/scripts"
	"time"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report [flags] [DAY] [TIMEZONE]",
	Short: "Generate a report on a specific day (YYYY-MM-DD) for a specific timezone",
	Long: `This command generates a report from the data present in the data directory based on the given day and timezone.
Be sure to use YYYY-MM-DD for the date and Area/Location according to the IANA timezone database for the timezone.`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			if len(args) == 0 {
				fmt.Println("No arguments provided. Check memento report --help for more information.")
			} else if len(args) == 1 {
				fmt.Println("Not enough arguments provided. Check memento report --help for more information")
			} else if len(args) > 2 {
				fmt.Println("Too many arguments provided. Check memento report --help for more information.")
			}
			return
		} else if _, ok := time.Parse("2006-01-02", args[0]); ok != nil {
			fmt.Println("Couldn't parse date - be sure to follow the format YYYY-MM-DD")
			return
		} else if _, ok := time.LoadLocation(args[1]); ok != nil {
			fmt.Println("Couldn't pares timezone - be sure to use Area/Location according to the IANA timezone database")
		} else {
			fileName, _ := cmd.Flags().GetString("file-name")

			timezone, _ := time.LoadLocation(args[1])
			utcDay, _ := time.Parse("2006-01-02 ", args[0])

			scripts.Reporter(utcDay, fileName, timezone)
		}
	},
}

func init() {
	RootCmd.AddCommand(reportCmd)

	reportCmd.Flags().String("file-name", "", "The filename of the report")
}
