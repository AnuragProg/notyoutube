package main

import (
	// "fmt"
	"time"
)

type Renderer struct {
	fps       uint16
	frameChan <-chan Frame
	terminal  Terminal
}

func NewRenderer(fps uint16, frameChan <-chan Frame) Renderer {
	return Renderer{
		fps:       fps,
		frameChan: frameChan,
		terminal:  NewTerminal(),
	}
}

func (r *Renderer) renderFrame(frame Frame) {
	frameWidth, frameHeight := frame.GetSize()

	for row := range frameHeight {
		for col := range frameWidth {
			rgb := frame.At(row, col)
			r.terminal.BatchPaintPixel(rgb.R, rgb.G, rgb.B)
		}
		r.terminal.BatchString("\n")
	}
	r.terminal.Flush()
}

func (r *Renderer) renderDeltaFrame(deltaFrame DeltaFrame) {
	for _, deltaPixel := range deltaFrame {
		r.terminal.BatchPaintPixelAt(uint32(deltaPixel.Row+1), uint32(deltaPixel.Col+1), deltaPixel.Pixel.R, deltaPixel.Pixel.G, deltaPixel.Pixel.B)
	}
	r.terminal.Flush()
}

func (r *Renderer) Render() error {
	// clear screen before anything
	r.terminal.BatchClearScreen()
	r.terminal.Flush()

	pause := time.Millisecond * time.Duration((float64(1) / float64(r.fps))*100)

	terminalWidth, terminalHeight, _ := r.terminal.GetSize()

	currentFrame := <-r.frameChan
	currentFrame.Resize(terminalWidth, terminalHeight)
	r.renderFrame(currentFrame)
	for nextFrame := range r.frameChan {
		nextFrame.Resize(terminalWidth, terminalHeight)
		deltaFrame, err := ComputeDeltaFrame(currentFrame, nextFrame)
		if err != nil {
			return err
		}
		time.Sleep(pause)
		r.renderDeltaFrame(deltaFrame)
	}
	return nil
}
