package scripts

import (
	"fmt"
	"log"
	"memento/structs"
	"memento/utils"
	"os"
	"time"
)

func Reporter(startDate time.Time, fileName string, timezone *time.Location) {
	// convert start and end date to the proper timezone
	startDate = startDate.In(timezone)
	endDate := startDate.AddDate(0, 0, 1)

	// load the data type enums here
	typeEnums := utils.LoadConfig()

	// make the report file
	if fileName == "" {
		fileName = startDate.Format("Report-2006-01-02")
	}
	file, err := os.Create("./reports/" + fileName + ".txt")
	if err != nil {
		log.Println(err)
	}

	// write the header of the report
	_, err = file.WriteString("Report " + startDate.Format("2006-01-02"))
	if err != nil {
		fmt.Println(err)
	}

	// read the required gob file
	var monthData structs.MonthData
	err = utils.ReadGob("./data/"+startDate.Format("2006-01")+".gob", &monthData)

	fmt.Println(endDate)
	// loop over the gob file and write to the report on every match
	for _, value := range monthData.Data {
		valueStart := value.Start.In(timezone)
		valueEnd := value.Start.In(timezone)

		// if the value is inside the desired timespan
		if valueStart.After(startDate) && valueStart.Before(endDate) {
			fileStat, _ := os.Stat(value.Path)
			modTime := fileStat.ModTime()

			fmt.Println(valueStart, modTime)

			if value.Dated == "point" {
				_, err := file.WriteString("\n" + valueStart.Format("15:04:05") + "            | " + typeEnums[value.Type].TrueName + "\t" + value.Path)
				if err != nil {
					log.Println(err)
				}
			} else if value.Dated == "range" {
				_, err := file.WriteString("\n" + valueStart.Format("15:04:05") + " - " + valueEnd.Format("15:04:05") + " | " + typeEnums[value.Type].TrueName + "\t\t" + value.Path)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
