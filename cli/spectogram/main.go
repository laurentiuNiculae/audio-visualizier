package main

import (
	"fmt"
	"math/cmplx"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/laurentiuNiculae/audio-visualizier/pkg/utils"
	"github.com/mjibson/go-dsp/fft"
)

const (
	width  = 1200
	height = 600

	graphStartX = 20
	graphStartY = int32(height / 2)
	graphEndX   = width - 20
	graphWidth  = graphEndX - graphStartX
	graphHeight = height - 100
	speed       = 2 * graphEndX
)

func main() {
	rl.InitWindow(width, height, "raylib waveform")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	if len(os.Args) == 1 {
		panic("give at least 1 .wav file name")
	}

	backgroundColor := rl.GetColor(0x181818FF)

	_, fData64 := utils.InitializeAudioData64(os.Args[1])

	hopLength := 512
	sampleWindow := 2048

	buff := make([][]float64, 0, graphWidth)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(backgroundColor)
		rl.DrawText("Press 'q' to quit or 'r' to restart.", 10, 15, 15, rl.White)

		if len(buff) != graphWidth {
			xCoordProgress := 0 // measured in pixels
			for i := 0; i < graphWidth && i*hopLength < len(fData64); i++ {
				fftFrameStart := i * hopLength
				fftFrameEnd := fftFrameStart + sampleWindow

				// get the component frequencies
				fftRez := fft.FFTReal(fData64[fftFrameStart:fftFrameEnd])

				// convert to real numbers by taking the abs, modulus
				freqMagnitutes := make([]float64, len(fftRez))

				for i := 0; i < len(fftRez); i++ {
					freqMagnitutes[i] = cmplx.Abs(fftRez[i])
				}

				// downsample the frequency magnitude slice to fit int the screen
				drawingBuffer := utils.Downsample(freqMagnitutes, graphHeight, utils.Max[float64])

				buff = append(buff, drawingBuffer)

				xCoordProgress += hopLength
			}
		}

		prog := float64(width) / float64(rl.GetMouseX())
		brightness := 1 + 1*prog

		fmt.Println(prog, brightness)

		for x := graphStartX; x < graphEndX && x < len(buff); x++ {
			for y := graphStartY - graphHeight/2; y < graphHeight; y++ {
				rl.DrawPixel(int32(x), int32(y), getColorFrom(buff[x][y], brightness))
			}
		}

		rl.EndDrawing()
	}
}

func getColorFrom(f float64, iso float64) rl.Color {
	brighness := uint8(utils.Min(f/iso*255, 255))

	return rl.NewColor(brighness, brighness, brighness, 255)
}
