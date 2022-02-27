package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"memento/structs"
	"os"
	"path"
	"path/filepath"
)

func LoadConfig(projectPath string) (map[string]structs.TypeEnum, error) {
	// load the data type enums here
	jsonFile, err := os.Open(filepath.Join(filepath.Clean(projectPath), `\config\typesEnum.json`))
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	var typeEnums map[string]structs.TypeEnum

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &typeEnums)
	if err != nil {
		return nil, err
	}

	return typeEnums, nil
}

func IsType(name, projectPath string) (bool, error) {
	config, err := LoadConfig(projectPath)
	if Handle(err) != nil {
		return false, err
	}
	if _, ok := config[name]; ok {
		return true, nil
	}
	return false, nil
}

func EnsureAppdata() (err error) {
	appDataPath := os.Getenv("APPDATA")
	configPath := "\\memento\\permSources.json"
	fullPath := path.Join(appDataPath, configPath)

	if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {
		// the appdata config file does not exist, it has to be created
		err = os.Mkdir(path.Join(appDataPath, "\\memento"), os.ModePerm)
		if err != nil {
			return err
		}

		err = os.WriteFile(fullPath, []byte("{}"), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadAppdata() (Appdata map[string][]structs.PermSource, err error) {
	appDataPath := os.Getenv("APPDATA")
	configPath := "\\memento\\permSources.json"
	fullPath := path.Join(appDataPath, configPath)

	jsonFile, err := os.Open(fullPath)
	if err != nil {
		return
	}

	defer jsonFile.Close()
	var permSources map[string][]structs.PermSource

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &permSources)
	if err != nil {
		return
	}

	return permSources, nil
}
