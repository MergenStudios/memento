package matcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
)


func GenerateBasePattern(path string) map[string]interface{} {
	var pattern = make(map[string]interface{})

	pattern["FilePath"] = path
	return pattern
}


func GetExifData(path string) (interface{}, error) {
	fmt.Println("called")
	cmd := exec.Command("exiftool", "-json", path)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var exifData []map[string]interface{}
	err = json.Unmarshal(out.Bytes(), &exifData)
	if err != nil {
		return nil, err
	}

	return exifData[0], nil
}

func GetFileName(path string) (interface{}, error) {
	basename := filepath.Base(path)
	ext := filepath.Ext(path)
	return filepath.Base(path)[:len(basename) - len(ext)], nil
}

func GetFileExtention(path string) (interface{}, error) {
	return filepath.Ext(path), nil
}

// This function is spaghetti code
// I am embracing the spagetti code
// can I please just have this in go:

// def a(): pass
// def b(): pass
// def c(): pass
// thing = {"a": a, "b": b, "c": c}
// thing["a"]()

// update: it works. in go. BUT ITS THE WORST THING I HAVE EVER LAYED MY EYES ON KILL ME NOW
// p.s.: dont acc kill me I wannna finish this project

func AddFilePattern(pattern map[string]interface{}, path, part string) error {
	var LookupTable map[string]func(string) (interface{}, error)
	LookupTable = map[string]func(string) (interface{}, error) {
		"ExifData": GetExifData,
		"FileName": GetFileName,
		"FileExtention": GetFileExtention,
	}

	if _, ok := LookupTable[part]; !ok {
		panic("panic") // TODO: make custom errors here
	}

	data, err := LookupTable[part](path)
	if err != nil {
		return err
	}

	pattern[part] = data
	return nil
}