package matcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func LoadPatterns(projectPath string) ([]map[string]interface{}, error) {
	jsonFile, err := os.Open(filepath.Join(filepath.Clean(projectPath), `\patterns.json`))
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	var patterns []map[string]interface{}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &patterns)
	if err != nil {
		return nil, err
	}


	fmt.Println(patterns)
	return patterns, nil
}
