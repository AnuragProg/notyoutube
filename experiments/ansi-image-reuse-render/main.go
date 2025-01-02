package main

import (
	"fmt"
	"os"
	"time"
	"golang.org/x/term"
)


func main(){

	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	fmt.Printf("Terminal size widthxheight = %vx%v\n", width, height)

	frameChan := make(chan Frame)
	defer close(frameChan)

	renderer := NewRenderer(60, frameChan)

	go func(){
		defer close(frameChan)
		for i:=1; i<=739; i++ {
			filename := fmt.Sprintf("./assets/frame_%03d.png", i)
			file, err := os.Open(filename)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			frame, err := NewFrame(file)
			if err != nil {
				panic(err)
			}
			frameChan<- frame
		}
		time.Sleep(time.Second*10)
	}()

	renderer.Render()
}
