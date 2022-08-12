package scripts

import (
	"memento/structs"
	"memento/utils"
	"os"
	"strings"
	"time"
)

func Reporter(startDate time.Time, fileName string, timezone *time.Location) {
	convertedDay := startDate.In(timezone)
	name, _ := convertedDay.Zone()

	startDate, _ = time.Parse("2006-01-02!MST", startDate.Format("2006-01-02")+"!"+name)
	endDate := startDate.AddDate(0, 0, 1)

	// make the report file
	var file *os.File
	if fileName == "" {
		fileName = "./reports/" + startDate.Format("Report-2006-01-02") + ".txt"
	} else if !strings.HasSuffix(fileName, ".txt") {
		fileName += ".txt"
	}
	file, err := os.Create(fileName)
	if utils.Handle(err) != nil {
		return
	}

	var reportLines []string

	// create the header of the report
	headerLine := "Report " + startDate.Format("2006-01-02") + "\n"
	var bodyLines []string

	// read the required gob file
	var monthData structs.MonthData
	err = utils.ReadGob("./data/"+startDate.Format("2006-01")+".gob", &monthData)
	if utils.Handle(err) != nil {
		return
	}

	var fileCount int64
	for _, value := range monthData.Data {
		valueStart := value.Start.In(timezone)

		// if the value is inside the desired timespan
		if valueStart.After(startDate) && valueStart.Before(endDate) {
			bodyLines = append(bodyLines, valueStart.Format("15:04:05")+"\t"+value.Path)

			fileCount++
		}
	}

	// TODO: make this return errors so it returns "no files found" from the command
	if fileCount == 0 {
		file.WriteString("No Files found for " + startDate.Format("2006-01-02"))
	}

	// write the file
	reportLines = append(reportLines, headerLine)
	reportLines = append(reportLines, bodyLines...)

	for _, line := range reportLines {
		_, err := file.WriteString(line + "\n")
		if utils.Handle(err) != nil {
			return
		}
	}
}
