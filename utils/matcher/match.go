package matcher

import (
	"reflect"
	"regexp"
)


// If I can *somehow* avoid it, I will never touch this function again
// its the best I could come up with:
// This function checks if the keys of one map are all contained in
// another map. If the value of a key is another interface, the functions
// recurses and calls itself.

func MatchPattern(configPattern map[string]interface{}, filePattern map[string]interface{}) (bool, error) {
	for configKey, configValue := range configPattern {
		if _, ok := filePattern[configKey]; !ok {
			err := AddFilePattern(filePattern, filePattern["FilePath"].(string), configKey)
			if err != nil {
				return false, err
			}
		}
		if fileValue, ok := filePattern[configKey]; ok {
			switch configValue.(type) {
			case string:
				stringConfigValue := configValue.(string)
				stringFileValue := fileValue.(string)
				
				if stringConfigValue == "+" {
					continue
				}

				regex := regexp.MustCompile(stringConfigValue)
				match := regex.MatchString(stringFileValue)

				if !match {
					return false, nil
				}
			case map[string]interface{}:
				mapFileValue := fileValue.(map[string]interface{})
				mapConfigValue := configValue.(map[string]interface{})

				check, err := MatchPattern(mapConfigValue, mapFileValue)
				if !check {
					return false, err
				}
			default:
				if reflect.DeepEqual(configValue, fileValue) {
					continue
				} else {
					return false, nil
				}
			}
		} else {
			return false, nil
		}
	}
	return true, nil
}

func MatchPatterns(patterns []map[string]interface{}, filePattern map[string]interface{}) (bool, map[string]interface{}, error) {
	for _, pattern := range patterns {
		cleanedPattern := make(map[string]interface{})
		for key, value := range pattern {
			cleanedPattern[key] = value
		}

		delete(cleanedPattern, "pattern")
		check, err := MatchPattern(cleanedPattern, filePattern)
		if err != nil {
			return false, nil, err
		}
		if check {
			return true, pattern, nil
		}
	}
	return false, nil, nil
}
