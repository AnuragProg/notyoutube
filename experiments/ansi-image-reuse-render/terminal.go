package main

import (
	"bufio"
	"fmt"
	"os"
	"golang.org/x/term"
)

// will be used to communicate with terminal
// will be structured to batch the output so that we don't have to face io latency issue
type Terminal struct {
	buffer *bufio.Writer
}

func NewTerminal() Terminal {
	return Terminal{
		buffer: bufio.NewWriter(os.Stdout),
	}
}

func (t *Terminal) GetSize()(width int, height int, err error){
	return term.GetSize(int(os.Stdout.Fd()))
}

func (t *Terminal) BatchString(str string) error {
	_, err := t.buffer.Write([]byte(str))
	return err
}

func (t *Terminal) BatchPaintPixel(r, g, b uint8) error {
	_, err := t.buffer.Write([]byte(fmt.Sprintf("\033[48;2;%v;%v;%vm ", r, g, b))) // here it should take 1 character because we are printing space
	return err
}

func (t *Terminal) BatchPaintPixelAt(x, y uint32, r, g, b uint8) error {
	t.BatchMoveCursorTo(x, y)
	_, err := t.buffer.Write([]byte(fmt.Sprintf("\033[48;2;%v;%v;%vm ", r, g, b))) // here it should take 1 character because we are printing space
	return err
}

func (t *Terminal) PaintPixelAt(x, y uint32, r, g, b uint8) error {
	t.BatchMoveCursorTo(x, y)
	_, err := t.buffer.Write([]byte(fmt.Sprintf("\033[48;2;%v;%v;%vm ", r, g, b))) // here it should take 1 character because we are printing space
	if err != nil {
		return err
	}
	t.buffer.Flush()
	return nil
}

func (t *Terminal) BatchMoveCursorTo(x, y uint32) error {
	// fmt.Println("moving cursor to x,y", x,y)
	_, err := t.buffer.Write([]byte(fmt.Sprintf("\033[%v;%vH", x, y)))
	return err
}

func (t *Terminal) Flush() error {
	return t.buffer.Flush()
}

func (t *Terminal) BatchClearScreen() error {
	_, err := t.buffer.Write([]byte("\033[2J"))
	return err
}
