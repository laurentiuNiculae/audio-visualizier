package main

import (
	"math/cmplx"
	"os"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/laurentiuNiculae/audio-visualizier/pkg/utils"
	"github.com/mjibson/go-dsp/fft"
)

const (
	width  = 1200
	height = 800
)

var (
	graphStartX = 20
	graphStartY = int32(height / 2)
	graphEndX   = width - 20
	graphWidth  = graphEndX - graphStartX
	graphHeight = height - 100
)

func main() {
	rl.InitWindow(width, height, "raylib waveform")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	if len(os.Args) == 1 {
		panic("give at least 1 .wav file name")
	}

	backgroundColor := rl.GetColor(0x181818FF)
	rl.ClearBackground(backgroundColor)

	_, fData64 := utils.InitializeAudioData64(os.Args[1])

	hopLength := 512 * 2
	sampleWindow := 2048 * 6

	graphWidth = utils.Min(graphWidth, len(fData64))

	buff := make([][]float64, 0, graphWidth)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.DrawText("Press 'q' to quit or 'r' to restart.", 10, 15, 15, rl.White)

		if len(buff) == 0 {
			xCoordProgress := 0 // measured in pixels
			window := make([]float64, sampleWindow)
			hammingWindow := utils.GetHammingValues(sampleWindow)

			for i := 0; i < graphWidth && i*hopLength < len(fData64); i++ {
				fftFrameStart := i * hopLength
				fftFrameEnd := utils.Min(fftFrameStart+sampleWindow, len(fData64)-1)

				// apply window
				copy(window, fData64[fftFrameStart:fftFrameEnd])
				for i := range window {
					window[i] = window[i] * hammingWindow[i]
				}

				// get the component frequencies
				fftRez := fft.FFTReal(window)

				// the top half is just symetric of the bottom half
				fftRez = fftRez[:len(fftRez)/2]

				// convert to real numbers by taking the abs, modulus
				freqMagnitutes := make([]float64, len(fftRez))

				for i := 0; i < len(fftRez); i++ {
					freqMagnitutes[i] = cmplx.Abs(fftRez[i])
				}

				// downsample the frequency magnitude slice to fit int the screen
				drawingBuffer := utils.LogDownsample(freqMagnitutes, graphHeight, utils.Max[float64])

				slices.Reverse(drawingBuffer)
				buff = append(buff, drawingBuffer)

				xCoordProgress += hopLength
			}

			for x := graphStartX; x < graphEndX && x < len(buff); x++ {
				for y := graphStartY - int32(graphHeight)/2; y < int32(graphHeight); y++ {
					rl.DrawPixel(int32(x), int32(y), getColorFrom(buff[x][y], 10))
				}
			}
		}

		rl.EndDrawing()
	}
}

func getColorFrom(f float64, iso float64) rl.Color {
	brighness := uint8(utils.Min(f/iso*255, 255))

	return rl.NewColor(brighness, brighness, brighness, 255)
}
