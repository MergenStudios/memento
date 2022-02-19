package utils

import (
	"encoding/gob"
	"os"
)

func WriteGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if Handle(err) != nil { return err}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	encoder.Encode(object)

	return nil
}

func ReadGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if Handle(err) != nil { return err}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(object)

	return nil
}
