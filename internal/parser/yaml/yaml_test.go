package yaml

import (
	"fmt"
	"testing"
)

func TestYaml(t *testing.T) {
    var yaml = `
stages:
  - prepare
  - build
  - test
  - deploy
    `
    result, err := Parse("config.yaml", yaml)
    if err != nil {
        panic("parser yaml failed")
    }
    fmt.Println(result["stages"])
}
