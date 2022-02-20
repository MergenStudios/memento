package scripts

import (
	"fmt"
	"memento/structs"
	"memento/utils"
	"os"
	"time"
)

func Reporter(startDate time.Time, fileName string, timezone *time.Location) {
	convertedDay := startDate.In(timezone)
	name, _ := convertedDay.Zone()

	startDate, _ = time.Parse("2006-01-02!MST", startDate.Format("2006-01-02")+"!"+name)
	endDate := startDate.AddDate(0, 0, 1)

	// load the data type enums here
	typeEnums, err := utils.LoadConfig()
	if utils.Handle(err) != nil {
		return
	}

	// make the report file
	if fileName == "" {
		fileName = startDate.Format("Report-2006-01-02")
	}
	file, err := os.Create("./reports/" + fileName + ".txt")
	if utils.Handle(err) != nil {
		return
	}

	// write the header of the report
	_, err = file.WriteString("Report " + startDate.Format("2006-01-02") + "\n")
	if utils.Handle(err) != nil {
		return
	}

	// read the required gob file
	var monthData structs.MonthData
	err = utils.ReadGob("./data/"+startDate.Format("2006-01")+".gob", &monthData)
	if utils.Handle(err) != nil {
		return
	}

	fmt.Println(startDate, endDate)
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
				if utils.Handle(err) != nil {
					return
				}
			} else if value.Dated == "range" {
				_, err := file.WriteString("\n" + valueStart.Format("15:04:05") + " - " + valueEnd.Format("15:04:05") + " | " + typeEnums[value.Type].TrueName + "\t\t" + value.Path)
				if utils.Handle(err) != nil {
					return
				}
			}
		}
	}
}
