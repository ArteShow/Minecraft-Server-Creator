package stats

import (
	"encoding/json"
	"fmt"
	"strings"
)

func SetValue(key string, value any, file []byte) ([]byte, error) {
	data := map[string]any{}
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	current := data
	parts := strings.Split(key, ".")
	for index, part := range parts[:len(parts)-1] {
		next, exists := current[part]
		if !exists || next == nil {
			child := map[string]any{}
			current[part] = child
			current = child
			continue
		}

		child, ok := next.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("key %q is not an object", strings.Join(parts[:index+1], "."))
		}

		current = child
	}

	current[parts[len(parts)-1]] = value

	newData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}

	return newData, nil
}