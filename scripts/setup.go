package scripts

import (
	"encoding/json"
	"memento/structs"
	"os"
)

func Setup() {
	// create the directory structure needed for a memento collection
	os.Mkdir("config", os.ModePerm)
	os.Mkdir("data", os.ModePerm)
	os.Mkdir("reports", os.ModePerm)

	typesEnum := map[string]structs.TypeEnum{
		"RECORDINGS": structs.TypeEnum{
			TrueName:      "Recording",
			Extensions:    []string{"mp4", "mov"},
			Dated:         "range",
			DetermineTime: "video",
		},
		"VIDEOS": structs.TypeEnum{
			TrueName:      "Video",
			Extensions:    []string{"mp4", "mov"},
			Dated:         "range",
			DetermineTime: "video",
		},
		"SCREENSHOTS": structs.TypeEnum{
			TrueName:      "Screenshot",
			Extensions:    []string{"png", "jpg"},
			Dated:         "point",
			DetermineTime: "mtime",
		},
		"PHOTOS": structs.TypeEnum{
			TrueName:      "Photo",
			Extensions:    []string{"png", "jpg"},
			Dated:         "point",
			DetermineTime: "mtime",
		},
	}

	jsonString, _ := json.Marshal(typesEnum)
	os.WriteFile("./config/typesEnum.json", jsonString, os.ModePerm)
}
