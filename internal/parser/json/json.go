package json

import (
	"encoding/json"
	"fmt"
)

func Parse(filename string, content string) (map[string]any, error) {
	result := make(map[string]any)
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		fmt.Printf("error parser %s: %s\n", filename, err)
		return nil, err
	}
	return result, nil
}
