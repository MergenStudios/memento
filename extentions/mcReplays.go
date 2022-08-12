package extentions

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

func McReplays(path string) (time.Time, error) {
	file, err := os.Open(path)
	if err != nil {
		return time.Time{}, err
	}
	defer file.Close()

	gzip_file, err := gzip.NewReader(file)
	if err != nil {
		return time.Time{}, err
	}
	defer gzip_file.Close()

	byteValues, _ := ioutil.ReadAll(gzip_file)

	// get the last timestamp from the log file
	check := regexp.MustCompile("\\[[0-9]{2}:[0-9]{2}:[0-9]{2}\\]")
	counter := len(byteValues) - 1
	line := ""
	for {
		if counter == 0 {
			line = string(byteValues)
			break
		}

		da_char := string(byteValues[counter])
		line = da_char + line

		if string(byteValues[counter]) == "\n" {
			if check.FindString(line) != "" {
				break
			}
		}
		counter--
	}

	// get the start time as a string from the beginning of the log file
	fmt.Println(line)
}