package main

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-audio/wav"
	"github.com/laurentiuNiculae/audio-visualizier/pkg/utils"
)

const (
	width  = 1200
	height = 600

	speed       = 5000
	graphStartX = 20
	graphStartY = int32(height / 2)
	graphEndX   = width - 20
	graphWidth  = graphEndX - graphStartX
	graphHeight = height - 200
)

type Selection struct {
	IsSelected bool
	Start      int32
	End        int32
}

type SampleOffsetRange struct {
	Down int32
	Up   int32
}

func main() {
	var (
		downscaleFrameWidth int32
		selection           Selection

		fData        []float32
		waveformBuff []SampleOffsetRange

		selectionColor  = rl.NewColor(30, 100, 100, 180)
		backgroundColor = rl.GetColor(0x181818FF)
	)

	rl.InitWindow(width, height, "raylib waveform")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	if len(os.Args) == 1 {
		panic("give at least 1 .wav file name")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	decod := wav.NewDecoder(file)
	decod.ReadInfo()

	PCM, err := decod.FullPCMBuffer()
	if err != nil {
		panic(err)
	}

	fData = PCM.AsFloat32Buffer().Data

	downscaleFrameWidth = utils.Max(1, int32(len(fData)/graphWidth))

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(backgroundColor)
		rl.DrawText("Press 'q' to quit or 'r' to restart.", 10, 15, 15, rl.White)

		gesture := rl.GetGestureDetected()

		if rl.IsKeyPressed(rl.KeyQ) {
			return
		} else if rl.IsKeyPressed(rl.KeyR) {
			rl.ClearBackground(backgroundColor)
			waveformBuff = []SampleOffsetRange{}

			if selection.IsSelected {
				dataLen := float64(len(fData))

				newStartIndex := uint64(dataLen * float64(selection.Start-graphStartX) / graphWidth)
				newEndIndex := uint64(dataLen * float64(selection.End-graphStartX) / graphWidth)

				if newStartIndex > newEndIndex {
					newStartIndex, newEndIndex = newEndIndex, newStartIndex
				}

				fmt.Println(newStartIndex, newEndIndex)
				fData = fData[newStartIndex:newEndIndex]
			} else {
				fData = PCM.AsFloat32Buffer().Data
			}

			downscaleFrameWidth = utils.Max(1, int32(len(fData)/graphWidth))
			selection.IsSelected = false
		} else if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if !selection.IsSelected {
				selection.IsSelected = true
			}

			selection.Start = utils.Clamp(rl.GetMouseX(), graphStartX, graphEndX)
			selection.End = selection.Start
		} else if gesture == rl.GestureDrag {
			selection.End = utils.Clamp(rl.GetMouseX(), graphStartX, graphEndX)
		} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			if selection.IsSelected {
				selection.End = utils.Clamp(rl.GetMouseX(), graphStartX, graphEndX)

				if utils.Abs(selection.End-selection.Start) < 5 {
					selection.IsSelected = false
					selection.Start = 0
					selection.End = 0
				}
			}
		}

		step := int32(len(waveformBuff))

		for i := 0; i < speed; i++ {
			if graphStartX+step < graphEndX {
				rl.DrawLine(20, height/2, graphStartX, height/2, rl.White)

				frame_start := downscaleFrameWidth * (step)
				frame_end := frame_start + downscaleFrameWidth

				min := utils.Min(fData[frame_start:frame_end]...)
				max := utils.Max(fData[frame_start:frame_end]...)

				upOffset := int32((graphHeight / 2) * utils.Abs(max))
				downOffset := int32((graphHeight / 2) * utils.Abs(min))

				waveformBuff = append(waveformBuff, SampleOffsetRange{
					Down: downOffset, Up: upOffset,
				})
			} else {
				break
			}

			step += 1
		}

		for i, offset := range waveformBuff {
			rl.DrawLine(graphStartX+int32(i), graphStartY+offset.Down, graphStartX+int32(i), graphStartY-offset.Up, rl.White)
		}

		rl.DrawLine(graphStartX, graphStartY, graphStartX+step, graphStartY, rl.White)

		if selection.IsSelected {
			selectionX := utils.Min(selection.Start, selection.End)
			selectionWidth := utils.Abs(selection.Start - selection.End)

			rl.DrawRectangle(selectionX, height/2-graphHeight/2, selectionWidth, graphHeight, selectionColor)
		}

		rl.EndDrawing()
	}
}
