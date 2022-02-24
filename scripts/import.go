package scripts

import (
	"fmt"
	"github.com/cnf/structhash"
	"github.com/superhawk610/bar"
	"memento/structs"
	"memento/utils"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ImportDatapoints(dataType string, inputPath string) {
	// set the timezone
	timezone, err := time.LoadLocation("UTC")
	if utils.Handle(err) != nil {
		return
	}

	// load the data type enums here
	typeEnums, err := utils.LoadConfig()
	if utils.Handle(err) != nil {
		return
	}
	dataPoints := make(map[string][]structs.DataPoint)

	fileCount, err := utils.FileCount(inputPath)
	if utils.Handle(err) != nil {
		return
	}

	bar1 := bar.NewWithOpts(
		bar.WithDimensions(fileCount, 50),
		bar.WithDisplay("[", "█", "█", " ", "]"),
		bar.WithFormat("Importing files :bar :percent"),
	)
	// make a list of all datapoints
	if _, err := os.Stat(inputPath); err == nil {
		// walk through every file in the path
		err = filepath.Walk(inputPath, func(filePath string, info os.FileInfo, err error) error {
			// fmt.Println(info.Name())

			// for every file extensions the give format can possibly have
			for _, extension := range typeEnums[dataType].Extensions {
				extension = "." + extension
				// if the extension of the file matches one of the possible extensions
				if extension == filepath.Ext(filePath) {
					// get the creation time of the file
					var startTime time.Time = time.Time{}
					var endTime time.Time = time.Time{}

					if typeEnums[dataType].DetermineTime == "mtime" {
						startTime = info.ModTime().In(timezone)
					} else if typeEnums[dataType].DetermineTime == "video" {

						ioReader, err := os.Open(filePath)
						if err != nil {
							return err
						}

						defer ioReader.Close()
						fileDuration, err := utils.GetMP4Duration(ioReader)
						if err != nil {
							return err
						}

						// calculate the start_time_parsed
						if fileDuration == 0 {
							startTime = info.ModTime()
							endTime = startTime
						}
						startTime = info.ModTime().Add(time.Duration(-(int64(fileDuration))) * time.Second)
						endTime = info.ModTime()
					}

					// creat the datapoint
					var dataPoint structs.DataPoint
					if typeEnums[dataType].Dated == "point" {
						dataPoint = structs.DataPoint{
							typeEnums[dataType].Dated,
							startTime,
							time.Time{},
							dataType,
							filePath,
						}
					} else if typeEnums[dataType].Dated == "range" {
						dataPoint = structs.DataPoint{
							typeEnums[dataType].Dated,
							startTime,
							endTime,
							dataType,
							filePath,
						}
					}

					// if the key doesnt exist, creat it and make the value an empty list of data points
					if _, ok := dataPoints[startTime.Format("2006-01")]; !ok {
						dataPoints[startTime.Format("2006-01")] = []structs.DataPoint{}
					}
					dataPoints[startTime.Format("2006-01")] = append(dataPoints[startTime.Format("2006-01")], dataPoint)
				}
			}
			bar1.TickAndUpdate(bar.Context{
				bar.Ctx("file", info.Name()),
			})
			return err
		})
		if utils.Handle(err) != nil {
			return
		}
	}

	fmt.Print("\n")
	bar2 := bar.NewWithOpts(
		bar.WithDimensions(len(dataPoints), 50),
		bar.WithDisplay("[", "█", "█", " ", "]"),
		bar.WithFormat("Serializing data :bar :percent"),
	)
	// write all the datapoints to the gob files
	for key, daDataPoints := range dataPoints {

		var monthData structs.MonthData
		gobPath := "./data/" + key + ".gob"

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
					fmt.Println("its a new one", datapointHash, monthData.Hashes)
					monthData.Hashes = append(monthData.Hashes, datapointHash)
					monthData.Data = append(monthData.Data, daDataPoints...)
				} else {
					fmt.Println("Duplicate found: " + datapoint.Path)
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
	return
}
