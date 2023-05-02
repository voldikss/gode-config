package json

import (
	"fmt"
	"testing"
)

func TestJsonParser(t *testing.T) {
	json := `
{
    "name": "alan",
    "age": 25
}
    `
	result, err := Parse("config.json", json)
	if err != nil {
		panic(err)
	}
    for k, v := range result {
        fmt.Println(k, v)
    }
}
