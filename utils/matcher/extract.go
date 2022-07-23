package matcher

import "fmt"

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
			fmt.Println("the path is over but we havent reached a string yet")
			return "", nil
		}

		extractedString, err := ExtractFromPath(pathArr[1:], thing.(map[string]interface{}))
		if err != nil {
			return "", err
		}

		return extractedString, nil

	default:
		fmt.Println("default case of the type switch in the path extractor")
		return "", nil
	}
}
