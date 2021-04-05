package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "Bearer teste"

	parts := strings.Split(str, "Bearer")

	if len(parts) == 2 {
		fmt.Println(parts[1])
		return
	}

	fmt.Println("Deu treta")
}
