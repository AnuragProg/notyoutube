package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
)

type RGB struct {
	R, G, B uint8
}

type Frame struct {
	frame [][]RGB
}

func NewFrame(frameInput io.Reader) (Frame, error) {
	image, _, err := image.Decode(frameInput)
	if err != nil {
		return Frame{}, err
	}

	imageSize := image.Bounds().Size()
	// fmt.Println(imageSize.X, imageSize.Y)
	frame := make([][]RGB, imageSize.Y)
	for idx := range frame {
		frame[idx] = make([]RGB, imageSize.X)
	}

	for row:=0; row<imageSize.Y; row++{
		for col:=0; col<imageSize.X; col++{

			color := image.At(col, row)
			r, g, b, _ := color.RGBA()

			frame[row][col] = RGB{
				R: uint8(r>>8), // because rgba is <<8 shifted
				G: uint8(g>>8), 
				B: uint8(b>>8),
			}
		}
	}

	return Frame{
		frame: frame,
	}, nil
}

func (f *Frame) At(row, col int) RGB{
	return f.frame[row][col]
}

func (f *Frame) Resize(newWidth, newHeight int) {
	currentWidth, currentHeight := f.GetSize()
	avgBlockWidth := int(math.Ceil(float64(currentWidth)/float64(newWidth)))
	avgBlockHeight := int(math.Ceil(float64(currentHeight)/float64(newHeight)))

	newFrame := make([][]RGB, newHeight)
	for idx := range newFrame{
		newFrame[idx] = make([]RGB, newWidth)
	}

	for row:=0; row<currentHeight; row+=avgBlockHeight {
		for col:=0; col<currentWidth; col+=avgBlockWidth {

			var blockR, blockG, blockB uint32 = 0, 0, 0
			var total uint32 = 0

			for _row:=row; _row<min(row+avgBlockHeight, currentHeight); _row++ {
				for _col:=col; _col<min(col+avgBlockWidth, currentWidth); _col++ {

					rgb := f.frame[_row][_col]
					blockR += uint32(rgb.R)
					blockG += uint32(rgb.G)
					blockB += uint32(rgb.B)
					total++
				}
			}

			blockR /= total
			blockG /= total
			blockB /= total

			// dividing by avgBlockHeight & avgBlockWidth, to retrieve original index for e.g interval 20 starting from 0 0/20=0 20/20=1 40/20=2
			// fmt.Println("inserting into x,y", row/avgBlockHeight, col/avgBlockWidth)
			newFrame[row/avgBlockHeight][col/avgBlockWidth] = RGB{
				R: uint8(blockR),
				G: uint8(blockG),
				B: uint8(blockB),
			}
		}
	}

	f.frame = newFrame
}

func (f *Frame) GetSize()(width, height int) {
	height = len(f.frame)
	if height == 0 {
		width = 0
		return
	}
	width = len(f.frame[0])
	return
}
