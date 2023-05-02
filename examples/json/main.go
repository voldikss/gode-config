package main

import (
	"fmt"
	"github.com/voldikss/gode-config/pkg/config"
)

func main() {
	// fmt.Println("examples json")

	name, err := config.Get("app.name")
	if err == nil {
		fmt.Println(name)
	}

	var m any
	config.GetAs("app", &m)
	fmt.Println(m)
	// fmt.Println(reflect.TypeOf(name).String())
	// fmt.Println(len(name.([]interface{})))
}
