package matcher

import (
	"memento/errors"
	"memento/extentions"
	"strings"
	"time"
)


// If a pattern matches, here the propper date time will be extracted
// This will do the regex extracting and the calling of the extention functions

// all functions will take in the filepath of the file in question
// and return a datetime and an errors
func GetDatetime(path string, filePattern, matchedPattern map[string]interface{}) (time.Time, error) {
	patternSpecifics := matchedPattern["pattern"].(map[string]string)


	if _, ok := patternSpecifics["Extention"]; ok {
		// todo: define the mapping somewhere else
		var ExtentionsMapping map[string]func(string) (time.Time, error)
		ExtentionsMapping = map[string]func(string) (time.Time, error) {
			"McReplays": extentions.McReplays,
		}

		extentionName := patternSpecifics["Extention"]
		startTime, err := ExtentionsMapping[extentionName](path)
		if err != nil {
			return time.Time{}, err
		}
		return startTime, nil

	} else {
		path := patternSpecifics["Path"]
		pathArr := strings.Split(path, "/")

		extractedString, err := ExtractFromPath(pathArr, filePattern)
		if err != nil {
			return time.Time{}, err
		}

		startTime, err := time.Parse(patternSpecifics["Format"], extractedString)
		if err != nil {
			return time.Time{}, err
		}

		return startTime, nil
	}
}


// This function extracts a string from a nested map, using a list
// of keys which functions as a path



func ExtractFromPath(pathArr []string, pattern map[string]interface{}) (string, error) {

	pathCompontent := pathArr[0]
	thing := pattern[pathCompontent]

	switch thing.(type) {
	case string:
		return thing.(string), nil
	case map[string]interface{}:
		if len(pathArr) == 1 {
			return "", &errors.FailedExtractionError{}
		}

		extractedString, err := ExtractFromPath(pathArr[1:], thing.(map[string]interface{}))
		if err != nil {
			return "", err
		}

		return extractedString, nil

	default:
		return "", &errors.FailedExtractionError{}
	}
}
