package utils

import "encoding/json"

func PrettifyJSON(o interface{}) string {
	prettyJSON, _ := json.MarshalIndent(o, "", "  ")
	return string(prettyJSON)
}

// Verify if a list contains a given element
func Contains[T comparable](list []T, x T) bool {
	for _, item := range list {
		if item == x {
			return true
		}
	}
	return false
}
