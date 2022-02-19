package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"memento/structs"
	"os"
)

func LoadConfig() (map[string]structs.TypeEnum, error) {
	// load the data type enums here
	jsonFile, err := os.Open("./config/typesEnum.json")
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	var typeEnums map[string]structs.TypeEnum

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &typeEnums)
	if err != nil {
		log.Println(err)
	}

	return typeEnums, nil
}

func IsType(name string) (bool, error) {
	config, err := LoadConfig()
	if Handle(err) != nil {return false, err}
	if _, ok := config[name]; ok {
		return true, nil
	}
	return false, nil
}
