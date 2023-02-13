package config

import (
	"go-config/internal/parser"
	"os"
	"strconv"
)

type Config struct {
	filename string
}

// TODO: maybe
// type Gonfig interface {
// 	Get(key string, defaultValue interface{}) (interface{}, error)
// 	GetString(key string, defaultValue interface{}) (string, error)
// 	GetInt(key string, defaultValue interface{}) (int, error)
// 	GetFloat(key string, defaultValue interface{}) (float64, error)
// 	GetBool(key string, defaultValue interface{}) (bool, error)
// 	GetAs(key string, target interface{}) error
// }

func (config *Config) Get(key string) (value any, err error) {
	return "", nil
}

func Get() {

}

func (config *Config) GetString(key string) (string, error) {
	value, err := config.Get(key)
	if err != nil {
		return "", err
	}
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case string:
		return v, nil
	}
	return "", &parser.InvalidValueError{Key: key, Kind: "string|int|int64"}
}

func LoadConfigFiles(filepath string) ([]byte, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func parseFile(filename string) {

}

func parse() {

}
