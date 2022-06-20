package utils

import (
	"encoding/json"
	"io/ioutil"
	"memento/structs"
	"os"
	"path/filepath"
)

func LoadPatterns(projectPath string) (map[string]structs.Pattern, error) {
	jsonFile, err := os.Open(filepath.Join(filepath.Clean(projectPath), `\patterns.json`))
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	var patterns map[string]structs.Pattern

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &patterns)
	if err != nil {
		return nil, err
	}

	return patterns, nil
}
