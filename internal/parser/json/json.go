package json

import (
	"encoding/json"
	"go-config/internal/types"
)

func Parse(filename string, content string) (types.ConfigType, error) {
	var result types.ConfigType
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, err
	}
	return result, nil
}
