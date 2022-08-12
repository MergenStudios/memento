package scripts

import (
	"fmt"
	"github.com/superhawk610/bar"
	"memento/structs"
	"memento/utils"
	"memento/utils/matcher"
	"os"
	"path/filepath"
)

func ImportDatapoints(inputPath, projectPath string) {
	fmt.Println(projectPath)
	patterns, err := matcher.LoadPatterns(projectPath)
	if utils.Handle(err) != nil {
		return
	}
	var dataPoints []structs.DataPoint

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

			// TODO: clean duplicates
			basePattern := matcher.GenerateBasePattern(filePath)
			check, pattern, err := matcher.MatchPatterns(patterns, basePattern)
			if err != nil {
				return err
			}

			if !check {
				// what do you do if it doesnt match a pattern
				bar1.Tick()
				return nil
			} else {
				patternInfo := pattern["pattern"].(map[string]string)
				name := patternInfo["Name"]

				startTime, err:= matcher.GetDatetime(filePath, basePattern, pattern)
				if err != nil {
					return err
				}

				dataPoint := structs.DataPoint{
					Start:	startTime,
					Path:	filePath,
					Type:	name,
				}

				dataPoints = append(dataPoints, dataPoint)

				bar1.Tick()
				return nil
			}
		})
	}

	// database stuff comes here


}