package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/voldikss/gode-config/internal/parser"
)

type Config struct {
	dir  string
	data map[string]any
}

const (
	godeEnvKey      = "GODE_ENV"
	configDirEnvKey = "GODE_CONFIG_DIR"
	configDir       = "config"
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
	configDir := os.Getenv(configDirEnvKey)
	if len(configDir) == 0 {
		cwd, _ := os.Getwd()
		configDir = filepath.Join(cwd, "config")
	}
	config = NewConfig(configDir)
}

var (
	configFileExtnames  = []string{".yaml", ".json", ".toml"}
	configFileBasenames = []string{"default"}
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

	goEnv := os.Getenv(godeEnvKey)
	if len(goEnv) == 0 {
		goEnv = "development"
	}

	configFileBasenames = append(configFileBasenames, goEnv, "local", "local"+goEnv)

	allowedFiles := make(map[string]int)
	resolutionIndex := 1
	for _, baseName := range configFileBasenames {
		for _, extName := range configFileExtnames {
			allowedFiles[baseName+extName] = resolutionIndex
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

	// fmt.Println(config.data)
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
	// case float32:
	// case float64:
	// return strconv.FormatFloat()
	case bool:
		return strconv.FormatBool(v), nil
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
	switch v := value.(type) {
	case string:
		return strconv.Atoi(v)
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case bool:
		if v == true {
			return 1, nil
		} else {
			return 0, nil
		}
	}
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
	switch v := value.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	// case int, int8, int16, int32, int64:
	// 	return
	case bool:
		if v == true {
			return 1.0, nil
		} else {
			return 0.0, nil
		}
	}

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

	switch v := value.(type) {
	case string:
		return strconv.ParseBool(v)
	case int, int8, int16, int32, int64, float32, float64:
		return v == 1, nil
	}

	return false, &ValueTypeNotMatchError{key, value, "bool"}
}

func GetAs(key string, target any) error {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidTargetTypeError{t: rv.Kind().String()}
	}

	value, err := Get(key)
	if err != nil {
		return err
	}

	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = json.Unmarshal(valueBytes, target)
	if err != nil {
		return err
	}
	return nil
}
