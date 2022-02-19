package scripts

import (
	"fmt"
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
	type_enums, err := utils.LoadConfig()
	if utils.Handle(err) != nil { return }
	data_points := make(map[string][]structs.DataPoint)

	if _, err := os.Stat(inputPath); err == nil {
		// walk through every file in the path
		err = filepath.Walk(inputPath, func(file_path string, info os.FileInfo, err error) error {
			fmt.Println(info.Name())

			// for every file extensions the give format can possibly have
			for _, extension := range type_enums[dataType].Extensions {
				extension = "." + extension
				// if the extension of the file matches one of the possible extensions
				if extension == filepath.Ext(file_path) {
					// get the creation time of the file
					var start_time time.Time = time.Time{}
					var end_time time.Time = time.Time{}

					if type_enums[dataType].DetermineTime == "mtime" {
						start_time = info.ModTime().In(timezone)
					} else if type_enums[dataType].DetermineTime == "video" {

						io_reader, err := os.Open(file_path)
						if err != nil {
							return err
						}

						defer io_reader.Close()
						file_duration, err := utils.GetMP4Duration(io_reader)
						if err != nil { return err}

						// calculate the start_time_parsed
						start_time = info.ModTime().Add(time.Duration(-(int64(file_duration))) * time.Second)
						end_time = info.ModTime()
					}


					// creat the datapoint
					var data_point structs.DataPoint
					if type_enums[dataType].Dated == "point" {
						data_point = structs.DataPoint{
							type_enums[dataType].Dated,
							start_time,
							time.Time{},
							dataType,
							file_path,
						}
					} else if type_enums[dataType].Dated == "range" {
						data_point = structs.DataPoint{
							type_enums[dataType].Dated,
							start_time,
							end_time,
							dataType,
							file_path,
						}
					}

					// if the key doesnt exist, creat it and make the value an empty list of data points
					if _, ok := data_points[start_time.Format("2006-01")]; !ok {
						data_points[start_time.Format("2006-01")] = []structs.DataPoint{}
					}
					data_points[start_time.Format("2006-01")] = append(data_points[start_time.Format("2006-01")], data_point)
				}
			}
			return err
		})
		if utils.Handle(err) != nil {
			return
		}
	}

	for key, da_data_points := range data_points {

		var month_data structs.MonthData
		gob_path := "./data/" + key + ".gob"

		// if the gob file exists
		if _, err := os.Stat(gob_path); err == nil {
			err = utils.ReadGob(gob_path, &month_data)
			if utils.Handle(err) != nil { return }

			month_data.Data = append(month_data.Data, da_data_points...)
		} else {
			split_string := strings.Split(key, "-")

			year, err := strconv.ParseInt(split_string[0], 10, 64)
			if utils.Handle(err) != nil { return }

			month, err := strconv.ParseInt(split_string[1], 10, 64)
			if utils.Handle(err) != nil { return }

			month_data.StartTime = time.Date(
				int(year),
				time.Month(month),
				00,
				00,
				00,
				00,
				000000000,
				timezone,
			)

			month_data.EndTime = month_data.StartTime.AddDate(0, 1, 0)

			month_data.Data = da_data_points
		}

		// sort the array of unsorted data points
		sort.Slice(month_data.Data[:], func(i, j int) bool {
			return month_data.Data[i].Start.Before(month_data.Data[j].Start)
		})

		// serialize the now sorted month data
		err := utils.WriteGob(gob_path, month_data)
		if utils.Handle(err) != nil { return }

	}
	return

	// TODO: add redundancy
	// TODO: solve issue when adding the same set of datapoints more than one time
	// TODO: figure out a way to handle errors
}
