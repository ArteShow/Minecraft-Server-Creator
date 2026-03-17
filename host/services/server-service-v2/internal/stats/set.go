package stats

import "encoding/json"

func SetValue(key string, value any, file []byte) []byte {
	data := map[string]any{}
	json.Unmarshal(file, &data)

	data[key] = value

	newData, _ := json.MarshalIndent(data, "", "  ")
	return newData
}