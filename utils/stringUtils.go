package utils

import "strings"

func CleanString(s string) string {
	newString := strings.Replace(s, "\n", "", -1)
	newString = strings.Replace(newString, "\r", "", -1)

	return newString
}
