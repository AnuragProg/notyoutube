package main

import (
	"fmt"
	"os"
	"time"
	// "golang.org/x/term"
)


func main(){
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
