package stats

import "encoding/json"

func GetValue(key string, file []byte) string {
	data := map[string]string{}
	json.Unmarshal(file, &data)

	return data[key]
}