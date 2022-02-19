package scripts

import (
	"encoding/json"
	"memento/structs"
	"memento/utils"
	"os"
)

func Setup() {
	// create the directory structure needed for a memento collection
	if utils.Handle(os.Mkdir("config", os.ModePerm)) != nil { return }
	if utils.Handle(os.Mkdir("data", os.ModePerm)) != nil { return }
	if utils.Handle(os.Mkdir("reports", os.ModePerm)) != nil { return }

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
	if utils.Handle(os.WriteFile("./config/typesEnum.json", jsonString, os.ModePerm)) != nil { return }
}
