package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"memento/structs"
	"os"
)

func LoadConfig() map[string]structs.TypeEnum {
	// load the data type enums here
	jsonFile, err := os.Open("./config/typesEnum.json")
	if err != nil {
		log.Println(err)
	}

	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Println(err)
		}
	}(jsonFile)
	var typeEnums map[string]structs.TypeEnum

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &typeEnums)
	if err != nil {
		log.Println(err)
	}

	return typeEnums
}

func IsType(name string) bool {
	config := LoadConfig()
	if _, ok := config[name]; ok {
		return true
	}
	return false
}
