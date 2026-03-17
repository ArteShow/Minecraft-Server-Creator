package stats

import "encoding/json"

func SetValue(key, value string, file []byte) []byte {
	data := map[string]string{}
	json.Unmarshal(file, &data)

	data[key] = value

	newData, _ := json.MarshalIndent(data, "", "  ")
	return newData
}