package main

import (
	"fmt"

	"on/internal/pkg"
)

func main() {
	port := pkg.Flags()
	fmt.Println(port)
}
