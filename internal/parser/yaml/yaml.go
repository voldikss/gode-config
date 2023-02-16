package yaml

import (
	"fmt"

	"github.com/go-yaml/yaml"
)

// TODO: yaml and yml use the same parser
func Parse(filename string, content string) (map[string]any, error) {
	result := make(map[string]any)
	err := yaml.Unmarshal([]byte(content), &result)
	if err != nil {
		return nil, err
	}
	fmt.Printf("--- m:\n%v\n\n", result)
	return result, nil
}
