package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func main(){
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	fmt.Printf("Terminal size widthxheight = %vx%v\n", width, height)
}
