package scripts

import (
	"fmt"
	"memento/utils"
)

func Update(path string) {
	permSources, err := utils.LoadAppdata()
	if utils.Handle(err) != nil {
		return
	}

	if path == "" {
		// update every path
		for projectPath, project := range permSources {
			for _, info := range project {
				ImportDatapoints(info.Type, info.Type, projectPath, false)
			}
		}
	} else {
		if _, ok := permSources[path]; !ok {
			fmt.Println("There are no permanent sources in this project yet, check memento import --help to find out how to add permanent sources.")
		} else {
			for _, info := range permSources[path] {
				ImportDatapoints(info.Type, info.Path, path, false)
			}
		}
	}
}
