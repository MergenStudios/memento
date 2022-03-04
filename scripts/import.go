package scripts

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cnf/structhash"
	"github.com/superhawk610/bar"
	"io"
	"memento/structs"
	"memento/utils"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ImportDatapoints(dataType, inputPath, projectPath string, permanent, console bool) {
	// set the timezone
	timezone, err := time.LoadLocation("UTC")
	if utils.Handle(err) != nil {
		return
	}

	// load the data type enums here
	typeEnums, err := utils.LoadConfig(projectPath)
	if utils.Handle(err) != nil {
		return
	}
	dataPoints := make(map[string][]structs.DataPoint)

	fileCount, err := utils.FileCount(inputPath)
	if utils.Handle(err) != nil {
		return
	}

	var bar1 *bar.Bar
	if console {
		bar1 = bar.NewWithOpts(
			bar.WithDimensions(fileCount, 50),
			bar.WithDisplay("[", "█", "█", " ", "]"),
			bar.WithFormat("Importing files  :bar :percent"),
		)
	}
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
			if console {
				bar1.Tick()
			}
			return err
		})
		if errors.Is(err, io.EOF) {
		} else if utils.Handle(err) != nil {
			return
		}
	}

	var bar2 *bar.Bar
	if console {
		fmt.Print("\n")
		bar2 = bar.NewWithOpts(
			bar.WithDimensions(len(dataPoints), 50),
			bar.WithDisplay("[", "█", "█", " ", "]"),
			bar.WithFormat("Serializing data :bar :percent"),
		)

	}
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

		if console {
			bar2.Tick()
		}
	}

	if console {
		fmt.Printf("\n")
		fmt.Println("Imported " +
			strconv.FormatInt(utils.GetDatapointsLen(dataPoints)-duplicates, 10) +
			" files, found " +
			strconv.FormatInt(duplicates, 10) +
			" duplicates out of " +
			strconv.FormatInt(utils.GetDatapointsLen(dataPoints), 10) +
			" files")
	}

	// add directory as a permanent data source
	if permanent {
		appDataPath := os.Getenv("APPDATA")
		configPath := "\\memento\\permSources.json"
		fullPath := path.Join(appDataPath, configPath)

		if utils.Handle(utils.EnsureAppdata()) != nil {
			return
		}

		permSources, err := utils.LoadAppdata()
		if utils.Handle(err) != nil {
			return
		}

		workingDir, err := utils.GetProjectPath()
		if err != nil {
			fmt.Println(err)
			return
		}

		if utils.Handle(err) != nil {
			return
		}

		permSources[workingDir] = append(permSources[workingDir], structs.PermSource{
			Type: dataType,
			Path: inputPath,
		})

		jsonString, err := json.Marshal(permSources)
		if utils.Handle(err) != nil {
			return
		}

		if utils.Handle(os.WriteFile(fullPath, jsonString, os.ModePerm)) != nil {
			return
		}
	}
	return
}
