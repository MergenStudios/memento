package matcher

import (
	"reflect"
	"regexp"
)


// If I can *somehow* avoid it, I will never touch this function again
// I know I absolutely suck at programming but its the best I could come up with:
// This function checks if the keys of one map are all contained in
// another map. If the value of a key is another interface, the functions
// recurses and calls itself.

func MatchPattern(configPattern map[string]interface{}, filePattern map[string]interface{}) bool {
	for configKey, configValue := range configPattern {
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
					return false
				}
			case map[string]interface{}:
				mapFileValue := fileValue.(map[string]interface{})
				mapConfigValue := configValue.(map[string]interface{})

				check := MatchPattern(mapConfigValue, mapFileValue)
				if !check {
					return false
				}
			default:
				if reflect.DeepEqual(configValue, fileValue) {
					continue
				} else {
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}

func MatchPatterns(patterns []map[string]interface{}, filePattern map[string]interface{}) (bool, map[string]interface{}) {
	for _, pattern := range patterns {
		claenedPattern := pattern
		delete(claenedPattern, "Use")

		check := MatchPattern(claenedPattern, filePattern)
		if check {
			return true, pattern
		}
	}
	return false, nil
}
