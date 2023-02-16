package config

import (
	"fmt"
	"go-config/internal/parser"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Config struct {
	dir  string
	data map[string]any
}

// TODO: maybe
// type Gonfig interface {
// 	Get(key string, defaultValue any) (any, error)
// 	GetString(key string, defaultValue any) (string, error)
// 	GetInt(key string, defaultValue any) (int, error)
// 	GetFloat(key string, defaultValue any) (float64, error)
// 	GetBool(key string, defaultValue any) (bool, error)
// 	GetAs(key string, target any) error
// }

const (
	GoEnvName        = "GO_ENV"
	ConfigDirEnvName = "GO_CONFIG_DIR"
	ConfigDirName    = "config"
)

var strictness = true

func GetStrictMode() bool {
	return strictness
}
func SetStrictMode(strict bool) {
	strictness = strict
}

var config *Config

func init() {
	configDir := os.Getenv(ConfigDirEnvName)
	if len(configDir) == 0 {
		cwd, _ := os.Getwd()
		configDir = filepath.Join(cwd, "config")
	}
	config = NewConfig(configDir)
}

// ///////////
var (
	ConfigFileExtNames  = []string{"yaml", "json", "toml"}
	ConfigFileBaseNames = []string{"default"}
)

func locateMatchedFiles(configDir string, allowedFiles map[string]int) []string {
	configFiles, _ := ioutil.ReadDir(ensureAbsolutePath(configDir))
	type FileWithResolutionIndex struct {
		index    int
		filename string
	}
	filesWithIndex := []FileWithResolutionIndex{}
	for _, configFile := range configFiles {
		configFileName := configFile.Name()
		if index, ok := allowedFiles[configFileName]; ok {
			filesWithIndex = append(filesWithIndex, FileWithResolutionIndex{index, configFileName})
		}
	}
	sort.Slice(filesWithIndex, func(i, j int) bool {
		return filesWithIndex[i].index < filesWithIndex[j].index
	})

	// TODO: improve
	result := []string{}
	for _, v := range filesWithIndex {
		result = append(result, filepath.Join(configDir, v.filename))
	}
	return result
}

func ensureAbsolutePath(p string) string {
	if p[0:1] == "." {
		cwd, _ := os.Getwd()
		return filepath.Join(cwd, p)
	}
	return p
}

func NewConfig(configDir string) *Config {
	data := make(map[string]any)

	goEnv := os.Getenv(GoEnvName)
	if len(goEnv) == 0 {
		goEnv = "development"
	}

	ConfigFileBaseNames = append(ConfigFileBaseNames, goEnv, "local", "local"+goEnv)

	allowedFiles := make(map[string]int)
	resolutionIndex := 1
	for _, baseName := range ConfigFileBaseNames {
		for _, extName := range ConfigFileExtNames {
			allowedFiles[baseName+"."+extName] = resolutionIndex
			resolutionIndex++
		}
	}

	config := &Config{configDir, data}

	configFiles := locateMatchedFiles(configDir, allowedFiles)
	for _, configFile := range configFiles {
		configObj, err := parser.ParseFile(configFile)
		if err == nil {
			config.extendDeep(configObj)
		}
	}

	// TODO: environment variables

	// TODO: cmd argument args

	fmt.Println(config.data)
	return config
}

func (config *Config) extendDeep(configObj map[string]any) {
	for k, v := range configObj {
		switch v.(type) {
		case map[string]any:
			target := config.data[k]
			if target == nil {
				config.data[k] = v
			} else {
				target := target.(map[string]any)
				source := v.(map[string]any)
				config.data[k] = extendDeepImpl(target, source)
			}
		default:
			config.data[k] = v
		}
	}
}

func extendDeepImpl(target, source map[string]any) map[string]any {
	for k, v := range source {
		switch v.(type) {
		case map[string]any:
			targetObj := target[k].(map[string]any)
			sourceObj := v.(map[string]any)
			target[k] = extendDeepImpl(targetObj, sourceObj)
		default:
			target[k] = v
		}
	}
	return target
}

func Has(property string) bool {
	return true
}

func Get(key string) (any, error) {
	fields := strings.Split(key, ".")
	value := getImpl(config.data, fields)
	if value == nil {
		return nil, &ValueNotFoundError{key}
	} else {
		return value, nil
	}
}

func getImpl(object map[string]any, fields []string) any {
	value := object[fields[0]]
	if len(fields) == 1 {
		return value
	}
	if obj, ok := value.(map[string]any); ok {
		return getImpl(obj, fields[1:])
	} else {
		return nil
	}
}

func GetString(key string) (string, error) {
	value, err := Get(key)
	if err != nil {
		return "", err
	}

	if stringValue, ok := value.(string); ok {
		return stringValue, nil
	}

	if GetStrictMode() {
		return "", &ValueTypeNotMatchError{key, value, "string"}
	}

	// TODO
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v), nil
	case float32, float64:
	}

	return "", &ValueTypeNotMatchError{key, value, "int"}
}

func GetInt(key string) (int, error) {
	value, err := Get(key)
	if err != nil {
		return -1, err
	}

	if intValue, ok := value.(int); ok {
		return intValue, nil
	}

	if GetStrictMode() {
		return -1, &ValueTypeNotMatchError{key, value, "int"}
	}

	// TODO
	// switch v := value.(type) {
	// case string:
	// }
	return -1, &ValueTypeNotMatchError{key, value, "int"}
}

func GetFloat(key string) (float64, error) {
	value, err := Get(key)
	if err != nil {
		return -1.0, err
	}

	if floatValue, ok := value.(float64); ok {
		return floatValue, nil
	}

	if GetStrictMode() {
		return -1.0, &ValueTypeNotMatchError{key, value, "float"}
	}

	// TODO
	// switch v.(type) {
	// case string:
	// }

	return -1.0, &ValueTypeNotMatchError{key, value, "float"}
}

func GetBool(key string) (bool, error) {
	value, err := Get(key)
	if err != nil {
		return false, err
	}

	if boolValue, ok := value.(bool); ok {
		return boolValue, nil
	}

	if GetStrictMode() {
		return false, &ValueTypeNotMatchError{key, value, "bool"}
	}
	// TODO
	// switch v.(type)

	return false, &ValueTypeNotMatchError{key, value, "bool"}
}

func GetAs(key string, target *any) error {
	value, err := Get(key)
	if err != nil {
		return err
	}
	*target = value
	return nil
}
