package utils

import "encoding/json"

func PrettifyJSON(o interface{}) string {
	prettyJSON, _ := json.MarshalIndent(o, "", "  ")
	return string(prettyJSON)
}
