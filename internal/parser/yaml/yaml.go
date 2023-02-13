package yaml

import (
	"fmt"

	"github.com/go-yaml/yaml"

	"go-config/internal/parser"
)

func Parse(filename string, content string) (map[any]any, error) {
	result := make(map[any]any)
	err := yaml.Unmarshal([]byte(content), &result)
	if err != nil {
		return nil, &parser.ParseError{Parser: "yaml"}
	}
	fmt.Printf("--- m:\n%v\n\n", result)
	return result, nil
}
