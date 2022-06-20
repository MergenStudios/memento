package scripts

import (
	"encoding/json"
	"memento/structs"
	"memento/utils"
	"os"
)

func Setup() {
	// create the directory structure needed for a memento collection
	if utils.Handle(os.Mkdir("data", os.ModePerm)) != nil {
		return
	}
	if utils.Handle(os.Mkdir("reports", os.ModePerm)) != nil {
		return
	}

	patterns := map[string]structs.Pattern{
		".png": {
			Regex:  "([0-9]{4}-[0-9]{2}-[0-9]{2}_[0-9]{2}[.][0-9]{2}[.][0-9]{2})",
			Format: "2006-01-02_15.04.05",
		},
		".gz": {
			Regex:  "([0-9]{4}-[0-9]{2}-[0-9]{2})-[0-9]*[.]log",
			Format: "2006.01.02",
		},
		".mp4": {
			Regex:  "([0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}-[0-9]{2}-[0-9]{2})",
			Format: "2006-01-02 15-04-05",
		},
		".mcpr": {
			Regex:  "([0-9]{4}_[0-9]{2}_[0-9]{2}_[0-9]{2}_[0-9]{2}_[0-9]{2})",
			Format: "2006_01_02_15_04_05",
		},
	}

	jsonString, _ := json.Marshal(patterns)
	if utils.Handle(os.WriteFile("./patterns.json", jsonString, os.ModePerm)) != nil {
		return
	}
}
