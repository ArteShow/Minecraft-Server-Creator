package stats

import (
	"encoding/json"
	"fmt"
	"strings"
)

func GetValue(key string, file []byte) (string, error) {
	data := map[string]any{}
	if err := json.Unmarshal(file, &data); err != nil {
		return "", err
	}

	current := any(data)
	for _, part := range strings.Split(key, ".") {
		object, ok := current.(map[string]any)
		if !ok {
			return "", fmt.Errorf("key %q is not an object path", key)
		}

		next, exists := object[part]
		if !exists {
			return "", fmt.Errorf("key %q not found", key)
		}

		current = next
	}

	switch value := current.(type) {
	case string:
		return value, nil
	case bool, float64, int, int64, nil:
		return fmt.Sprint(value), nil
	default:
		encoded, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		return string(encoded), nil
	}
}