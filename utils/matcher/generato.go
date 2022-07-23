package matcher

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
)

func GetExifData(path string) (map[string]interface{}, error) {
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

func GetFileName(path string) (string, error) {
	return filepath.Base(path), nil
}

func GetFileExtention(path string) (string, error) {
	return filepath.Ext(path), nil
}

func GenerateFilePattern(path string) (map[string]interface{}, error) {
	var filePattern = make(map[string]interface{})

	exifData, err := GetExifData(path)
	if err != nil {
		return nil, err
	}

	fileName, err := GetFileName(path)
	if err != nil {
		return nil, err
	}

	fileExtention, err := GetFileExtention(path)
	if err != nil {
		return nil, err
	}

	filePattern["ExifData"] = exifData
	filePattern["FileName"] = fileName
	filePattern["FileExtention"] = fileExtention

	return filePattern, nil
}