package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(width, height)

}
