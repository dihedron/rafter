package logging

import "encoding/json"

func ToJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func ToPrettyJSON(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "  ")
	return string(data)
}
