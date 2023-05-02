package main

import (
	"fmt"
	"reflect"

	"github.com/voldikss/gode-config/pkg/config"
)

type Config struct {
	App struct {
		name    string
		version string
	}
	Environment string

	List []struct {
		Name string
		Age  int
	}
}

func main() {
	// fmt.Println("examples yaml")

	list, err := config.Get("members.list")
	if err == nil {
		fmt.Printf("name: %v\n", list)
	}
	fmt.Println(reflect.TypeOf(list).String())
	for _, p := range list.([]any) {
		fmt.Println(p)
	}

	fmt.Println("------------------------------------")

	var c Config
	err = config.GetAs("members.list", &c.List)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("c.list: %v\n", c.List)

	// fmt.Println(name)
	// fmt.Println(reflect.TypeOf(name).String())
	// fmt.Println(len(name.([]interface{})))
}
