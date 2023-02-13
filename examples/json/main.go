package main

import (
	"fmt"
	"go-config/pkg/config"
)

func main() {
	fmt.Println("examples")

    name := config.Get()
}
