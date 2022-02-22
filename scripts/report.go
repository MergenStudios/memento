package scripts

import (
	"memento/structs"
	"memento/utils"
	"os"
	"strconv"
	"strings"
	"time"
)

func Reporter(startDate time.Time, fileName string, timezone *time.Location, stats bool) {
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
	var file *os.File
	if fileName == "" {
		fileName = "./reports/" + startDate.Format("Report-2006-01-02") + ".txt"
	} else if !strings.HasSuffix(fileName, ".txt") {
		fileName += ".txt"
	}
	file, err = os.Create(fileName)
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
	fileTypesCount := make(map[string]int64)
	for _, value := range monthData.Data {
		valueStart := value.Start.In(timezone)
		valueEnd := value.Start.In(timezone)

		// if the value is inside the desired timespan
		if valueStart.After(startDate) && valueStart.Before(endDate) {
			if value.Dated == "point" {
				bodyLines = append(bodyLines, valueStart.Format("15:04:05")+"            | "+typeEnums[value.Type].TrueName+"\t"+value.Path)

				fileCount++
				if _, ok := fileTypesCount[value.Type]; ok {
					fileTypesCount[value.Type]++
				} else {
					fileTypesCount[value.Type] = 1
				}
			} else if value.Dated == "range" {
				bodyLines = append(bodyLines,
					valueStart.Format("15:04:05")+
						" - "+
						valueEnd.Format("15:04:05")+
						" | "+
						typeEnums[value.Type].TrueName+
						"\t\t"+
						value.Path)

				fileCount++
				if _, ok := fileTypesCount[value.Type]; ok {
					fileTypesCount[value.Type]++
				} else {
					fileTypesCount[value.Type] = 1
				}
			}
		}
	}

	var statsLine string
	statsLine += strconv.FormatInt(fileCount, 10) + " files" + "\t"
	for key, val := range fileTypesCount {
		if len(fileTypesCount) == 1 {
			statsLine +=
				strconv.FormatInt(val, 10) +
					" x " +
					typeEnums[key].TrueName
		} else {
			statsLine +=
				strconv.FormatInt(val, 10) +
					" x " +
					typeEnums[key].TrueName +
					" ⁞ "
		}
	}
	statsLine += "\n"

	// write the file
	reportLines = append(reportLines, headerLine)
	if stats {
		reportLines = append(reportLines, statsLine)
	}
	reportLines = append(reportLines, bodyLines...)

	for _, line := range reportLines {
		_, err := file.WriteString(line + "\n")
		if utils.Handle(err) != nil {
			return
		}
	}
}
