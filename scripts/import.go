package scripts

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
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
	
	fmt.Println(fileCount)
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
			if filePath == inputPath {
				return nil
			}
			// TODO: clean duplicates
			basePattern := matcher.GenerateBasePattern(filePath)
			check, pattern, err := matcher.MatchPatterns(patterns, basePattern)
			if err != nil {
				return err
			}

			if !check {
				//dataPoint := structs.DataPoint{
				//	Start: time.Time{},
				//	Path: filePath,
				//	Type: "file",
				//}

				// dataPoints = append(dataPoints, dataPoint)
			} else {
				patternInfo := pattern["pattern"].(map[string]interface{})
				name := patternInfo["Name"].(string)

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
			}
			bar1.Tick()
			return nil
		})

	}


	var bar2 *bar.Bar
	bar2 = bar.NewWithOpts(
		bar.WithDimensions(len(dataPoints), 50),
		bar.WithDisplay("[", "█", "█", " ", "]"),
		bar.WithFormat("Adding files to database  :bar :percent"),
	)

	// database stuff comes here
	dbPath := projectPath + `\memento.db`
	db, err := sql.Open("sqlite3", dbPath)
	defer db.Close()
	if utils.Handle(err) != nil { return }


	querry := `
INSERT INTO DataPoints (start_time, file_path, type)
VALUES ($1, $2, $3)
`

	preppedQuery, err := db.Prepare(querry)
	if utils.Handle(err) != nil {
		return
	}

	for _, dataPoint := range dataPoints {
		startUnix := utils.PositiveTimestmap(dataPoint.Start.Unix())

		// todo: make this a bit more elegant
		if startUnix < 0 {
			_, err := preppedQuery.Exec(nil, dataPoint.Path, dataPoint.Type)
			if utils.Handle(err) != nil {
				return
			}
		} else {
			_, err := preppedQuery.Exec(startUnix, dataPoint.Path, dataPoint.Type)
			if utils.Handle(err) != nil {
				return
			}
		}

		bar2.Tick()
	}
}
