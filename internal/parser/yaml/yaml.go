package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

// TODO: yaml and yml use the same parser
func Parse(filename string, content string) (map[string]any, error) {
	result := make(map[string]any)
	if err := yaml.Unmarshal([]byte(content), &result); err != nil {
		fmt.Printf("error parser %s: %s\n", filename, err)
		return nil, err
	}
	return result, nil
}
