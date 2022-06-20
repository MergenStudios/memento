package scripts

import (
	"fmt"
	"github.com/cnf/structhash"
	"github.com/superhawk610/bar"
	"memento/structs"
	"memento/utils"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ImportDatapoints(inputPath, projectPath string) {
	// set the timezone
	timezone, err := time.LoadLocation("UTC")
	if utils.Handle(err) != nil {
		return
	}

	patterns, err := utils.LoadPatterns(projectPath)
	if utils.Handle(err) != nil {
		return
	}
	dataPoints := make(map[string][]structs.DataPoint)

	fileCount, err := utils.FileCount(inputPath)
	if utils.Handle(err) != nil {
		return
	}

	var bar1 *bar.Bar
	bar1 = bar.NewWithOpts(
		bar.WithDimensions(fileCount, 50),
		bar.WithDisplay("[", "█", "█", " ", "]"),
		bar.WithFormat("Importing files  :bar :percent"),
	)

	// make a list of all datapoints
	if _, err := os.Stat(inputPath); err == nil {
		// walk through every file in the path
		err = filepath.Walk(inputPath, func(filePath string, info os.FileInfo, err error) error {
			// fmt.Println(info.Name())
			fileExtention := filepath.Ext(filePath)
			if _, ok := patterns[fileExtention]; ok {
				pattern := patterns[fileExtention]
				regex := regexp.MustCompile(pattern.Regex)
				nameMatch := regex.FindStringSubmatch(info.Name())[0]

				startTime, err := time.Parse(pattern.Format, nameMatch)
				if utils.Handle(err) != nil {
					return err
				}

				dataPoint := structs.DataPoint{
					Start: startTime,
					Path:  filePath,
				}

				// if the key doesnt exist, creat it and make the value an empty list of data points
				if _, ok := dataPoints[startTime.Format("2006-01")]; !ok {
					dataPoints[startTime.Format("2006-01")] = []structs.DataPoint{}
				}
				dataPoints[startTime.Format("2006-01")] = append(dataPoints[startTime.Format("2006-01")], dataPoint)
			}
			bar1.Tick()
			return nil
		})
	}

	var bar2 *bar.Bar
	fmt.Print("\n")
	bar2 = bar.NewWithOpts(
		bar.WithDimensions(len(dataPoints), 50),
		bar.WithDisplay("[", "█", "█", " ", "]"),
		bar.WithFormat("Serializing data :bar :percent"),
	)

	// write all the datapoints to the gob files
	var duplicates int64 = 0
	for key, daDataPoints := range dataPoints {

		var monthData structs.MonthData
		gobPath := filepath.Join(projectPath, "data", key+".gob")

		// if the gob file exists
		if _, err := os.Stat(gobPath); err == nil {
			err = utils.ReadGob(gobPath, &monthData)
			if utils.Handle(err) != nil {
				return
			}
			for _, datapoint := range daDataPoints {
				datapointHash, err := structhash.Hash(datapoint, 1)
				if utils.Handle(err) != nil {
					return
				}

				if !utils.InList(datapointHash, monthData.Hashes) {
					monthData.Hashes = append(monthData.Hashes, datapointHash)
					monthData.Data = append(monthData.Data, daDataPoints...)
				} else {
					duplicates++
				}
			}

		} else {
			splitString := strings.Split(key, "-")

			year, err := strconv.ParseInt(splitString[0], 10, 64)
			if utils.Handle(err) != nil {
				return
			}

			month, err := strconv.ParseInt(splitString[1], 10, 64)
			if utils.Handle(err) != nil {
				return
			}

			// TODO: this can be improved
			monthData.StartTime = time.Date(
				int(year),
				time.Month(month),
				00,
				00,
				00,
				00,
				000000000,
				timezone,
			)

			monthData.EndTime = monthData.StartTime.AddDate(0, 1, 0)
			monthData.Data = daDataPoints

			for _, datapoint := range daDataPoints {
				datapointHash, err := structhash.Hash(datapoint, 1)
				if utils.Handle(err) != nil {
					return
				}
				monthData.Hashes = append(monthData.Hashes, datapointHash)
			}
		}

		// sort the array of unsorted data points
		sort.Slice(monthData.Data[:], func(i, j int) bool {
			return monthData.Data[i].Start.Before(monthData.Data[j].Start)
		})

		// serialize the now sorted month data
		err := utils.WriteGob(gobPath, monthData)
		if utils.Handle(err) != nil {
			return
		}

		bar2.Tick()
	}

	fmt.Printf("\n")
	fmt.Println("Imported " +
		strconv.FormatInt(utils.GetDatapointsLen(dataPoints)-duplicates, 10) +
		" files, found " +
		strconv.FormatInt(duplicates, 10) +
		" duplicates out of " +
		strconv.FormatInt(utils.GetDatapointsLen(dataPoints), 10) +
		" files")

	return
}
