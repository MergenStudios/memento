package utils

import (
	"errors"
	"os"
	"path/filepath"
)

func GetProjectPath() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	path := workingDir
	for {
		if path == "C:\\" {
			return "not found", nil
		}

		if _, err := os.Stat(filepath.Join(path, "patterns.json")); err == nil {
			return path, nil
		} else if errors.Is(err, os.ErrNotExist) {
			path = filepath.Clean(filepath.Join(path, ".."))
			continue
		} else {
			return "", err
		}
	}
}
