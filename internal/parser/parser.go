package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go-config/internal/parser/json"
	"go-config/internal/parser/yaml"
)

type Parser interface {
	Parse(filename string, content string) (any, error)
}

var ParserMap = map[string]func(string, string) (map[string]any, error){
	"json": func(filepath, content string) (map[string]any, error) {
		return json.Parse(filepath, content)
	},
	"yaml": func(filepath, content string) (map[string]any, error) {
		return yaml.Parse(filepath, content)
	},
}

func ParseFile(configFilePath string) (map[string]any, error) {
	data, _ := os.ReadFile(configFilePath)
	ext := filepath.Ext(configFilePath)[1:]
	parseFunc, ok := ParserMap[ext]
	if !ok || parseFunc == nil {
		fmt.Println("invalid config file with extension", ext, configFilePath)
		return nil, errors.New("invalid config file")
	}
	config, err := parseFunc(configFilePath, string(data))
	if err != nil {
		return nil, err
	} else {
		return config, nil
	}
}
