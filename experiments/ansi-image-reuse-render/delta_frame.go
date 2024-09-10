package main

import "errors"

type DeltaPixel struct {
	Row, Col int
	Pixel    RGB
}

type DeltaFrame []DeltaPixel

func ComputeDeltaFrame(fromFrame, toFrame Frame) (DeltaFrame, error) {
	deltaFrame := DeltaFrame{}

	frame1Width, frame1Height := fromFrame.GetSize()
	frame2Width, frame2Height := toFrame.GetSize()

	if frame1Width != frame2Width || frame1Height != frame2Height {
		return DeltaFrame{}, errors.New("frames must have same dimensions")
	}

	for row := 0; row < frame1Height; row++ {
		for col := 0; col < frame1Width; col++ {
			if fromFrame.At(row, col) != toFrame.At(row, col) {
				deltaFrame = append(deltaFrame, DeltaPixel{
					Row: row,
					Col: col,
					Pixel: toFrame.At(row, col),
				})
			}
		}
	}

	return deltaFrame, nil
}
