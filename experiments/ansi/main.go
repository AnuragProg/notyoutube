package main

import (
	"bufio"
	"time"

	// "flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"golang.org/x/term"
)

func CalculateAveragingBlockDimension(termWidth, termHeight, imgWidth, imgHeight int) (width int, height int) {
	width = int(math.Ceil(float64(imgWidth) / float64(termWidth)))
	height = int(math.Ceil(float64(imgHeight) / float64(termHeight)))
	return
}

func Render(img image.Image) {
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	imageDim := img.Bounds().Size()
	fmt.Printf("image widthxheight %vx%v\n", imageDim.X, imageDim.Y)

	avgBlockWidth, avgBlockHeight := CalculateAveragingBlockDimension(width, height, imageDim.X, imageDim.Y)

	fmt.Printf("block widthxheight %vx%v\n", avgBlockWidth, avgBlockHeight)
	screen := [][]uint32{}

	for y := 0; y < imageDim.Y; y += avgBlockHeight {
		row := []uint32{}
		for x := 0; x < imageDim.X; x += avgBlockWidth {

			var avgPixel uint32 = 0
			var total uint32 = 0

			for _y := y; _y < min(y+avgBlockHeight, imageDim.Y); _y += avgBlockHeight {
				for _x := x; _x < min(x+avgBlockWidth, imageDim.X); _x += avgBlockWidth {
					r, g, b, _ := img.At(_x, _y).RGBA()

					r = r >> 8
					g = g >> 8
					b = b >> 8

					pixelValue := ((r + g + b) / 3) // value between 0 - 254 inclusive

					// linear interpolation
					// convert x in range [a,b] to [c,d]
					// new_x = c + ((x-a)(d-c) / (b-a))
					// make it 16 - 231 ( for ansi to work in rgb mode)

					// don't know what kind of color
					// pixelValue = 16 + ((pixelValue-0)*(231-16))/254

					// grayscale
					// pixelValue = 232 + ((pixelValue-0)*(255-232))/254

					// 16 color
					// pixelValue = 0 + ((pixelValue-0)*(15-0))/254

					avgPixel += pixelValue
					total++
				}
			}
			row = append(row, avgPixel/total)
		}
		screen = append(screen, row)
	}

	fmt.Printf("averaging done with widthxheight %vx%v\n", len(screen[0]), len(screen))

	buffer := bufio.NewWriter(os.Stdout)
	for _, row := range screen {
		for _, pixel := range row {
			buffer.Write([]byte(fmt.Sprintf("\033[48;5;%dm ", pixel)))
			// fmt.Printf("\033[48;5;%dm ", pixel)
		}
		buffer.Write([]byte("\n"))
		// fmt.Println()
	}

	buffer.Flush()
	// for i := 40; i <= 47; i++ {
	// 	fmt.Printf("\033[48;5;%dm  COLOR %d  \033[0m\n", i, i)
	// }
}

func Clear() {
	fmt.Print("\033[2J")
}

func print[T any](arr []T) {
	fmt.Printf("Array addr %p: ", &arr)
	fmt.Print("members ")
	for idx := range arr {
		fmt.Printf("%p ", &arr[idx])
	}
	fmt.Println()
}

func TrueColorRender(img image.Image) {
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	imageDim := img.Bounds().Size()
	// fmt.Printf("image widthxheight %vx%v\n", imageDim.X, imageDim.Y)

	avgBlockWidth, avgBlockHeight := CalculateAveragingBlockDimension(width, height, imageDim.X, imageDim.Y)

	// fmt.Printf("block widthxheight %vx%v\n", avgBlockWidth, avgBlockHeight)
	screen := [][]string{}

	for y := 0; y < imageDim.Y; y += avgBlockHeight {
		row := []string{}
		for x := 0; x < imageDim.X; x += avgBlockWidth {

			var blockR, blockG, blockB uint32 = 0, 0, 0
			var total uint32 = 0

			for _y := y; _y < min(y+avgBlockHeight, imageDim.Y); _y++{
				for _x := x; _x < min(x+avgBlockWidth, imageDim.X); _x++{
					r, g, b, _ := img.At(_x, _y).RGBA()

					r = r >> 8
					g = g >> 8
					b = b >> 8

					blockR += r
					blockG += g
					blockB += b
					total++

				}
			}

			blockR /= total
			blockG /= total
			blockB /= total

			row = append(row, fmt.Sprintf("\033[48;2;%v;%v;%vm ", blockR, blockG, blockB))

		}
		screen = append(screen, row)
	}

	// fmt.Printf("averaging done with widthxheight %vx%v\n", len(screen[0]), len(screen))

	buffer := bufio.NewWriter(os.Stdout)
	for _, row := range screen {
		for _, pixel := range row {
			buffer.Write([]byte(pixel))
			// buffer.Write([]byte(fmt.Sprintf("\033[48;5;%dm ", pixel)))
			// fmt.Printf("\033[48;5;%dm ", pixel)
		}
		buffer.Write([]byte("\n"))
		// fmt.Println()
	}

	buffer.Flush()
	// for i := 40; i <= 47; i++ {
	// 	fmt.Printf("\033[48;5;%dm  COLOR %d  \033[0m\n", i, i)
	// }
}
func main() {

	for i := 1; i < 740; i++ {
		filename := fmt.Sprintf("./assets/frame_%03d.png", i)
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		img, _, err := image.Decode(file)
		if err != nil {
			panic(err)
		}
		TrueColorRender(img)
		// Render(img)
		time.Sleep(time.Millisecond * 10)
		Clear()
	}

	// var filename string
	// flag.StringVar(&filename, "f", "", "")
	// flag.Parse()
	// file, err := os.Open(filename)
	// if err != nil {
	// 	panic(err)
	// }
	// img, _, err := image.Decode(file)
	// if err != nil {
	// 	panic(err)
	// }
	// Render(img)
	// TrueColorRender(img)
	// Clear()
}
