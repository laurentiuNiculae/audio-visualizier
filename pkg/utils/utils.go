package utils

import (
	"math"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type numbers interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 | uint | uint8 | uint16 | uint32 | uint64
}

func Min[T numbers](vals ...T) T {
	if len(vals) == 0 {
		return 0
	}

	min := vals[0]

	for i := range vals {
		if vals[i] < min {
			min = vals[i]
		}
	}

	return min
}

func Max[T numbers](vals ...T) T {
	if len(vals) == 0 {
		return 0
	}

	max := vals[0]

	for i := range vals {
		if vals[i] > max {
			max = vals[i]
		}
	}

	return max
}

func Avg[T numbers](vals ...T) T {
	var avg T = 0

	for i := range vals {
		avg += vals[i]
	}

	return avg / T(len(vals))
}

func Abs[T numbers](x T) T {
	if x < 0 {
		return -x
	}

	return x
}

func Clamp[T numbers](x, low, high T) T {
	if x < low {
		return low
	}

	if x > high {
		return high
	}

	return x
}

func InitializeAudioData(file string) (*audio.IntBuffer, []float32) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	decod := wav.NewDecoder(f)
	decod.ReadInfo()

	PCM, err := decod.FullPCMBuffer()
	if err != nil {
		panic(err)
	}

	fData := PCM.AsFloat32Buffer().Data

	return PCM, fData
}

func InitializeAudioData64(file string) (*audio.IntBuffer, []float64) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	decod := wav.NewDecoder(f)
	decod.ReadInfo()

	PCM, err := decod.FullPCMBuffer()
	if err != nil {
		panic(err)
	}

	fData := PCM.AsFloat32Buffer().Data

	fData64 := make([]float64, len(fData))

	for i := range fData {
		fData64[i] = float64(fData[i])
	}

	return PCM, fData64
}

func Downsample[T numbers](slice []T, newSize int, reducer func(subSlice ...T) T) []T {
	result := make([]T, newSize)
	progress := float64(0.0)

	lastProcessedIndex := 0
	step := float64(len(slice)) / float64(newSize)

	// we'll iterate untill right before the last element
	// float addition can have errors so the last element might be lost if we're not careful
	for i := 0; i < newSize-1; i++ {
		progress += step

		poolingIntervalSize, _ := math.Modf(progress)

		start := lastProcessedIndex
		end := start + int(poolingIntervalSize)
		result[i] = reducer(slice[start:end]...)

		lastProcessedIndex = end
		progress -= poolingIntervalSize
	}

	result[newSize-1] = reducer(slice[lastProcessedIndex:]...)

	return result
}

func LogDownsample[T numbers](slice []T, newSize int, reducer func(subSlice ...T) T) []T {
	result := make([]T, newSize)
	lastProcessedIndex := 0
	factor := math.Pow(float64(len(slice)), float64(1)/float64(newSize))

	i := float64(factor)
	outIndex := 0
	progress := float64(0)
	poolingIntervalSize := float64(0)
	// we'll iterate untill right before the last element
	// float addition can have errors so the last element might be lost if we're not careful
	for ; i*factor < float64(len(slice)); i *= factor {
		step := i*factor - i
		progress += step

		poolingIntervalSize, progress = math.Modf(progress)

		start := lastProcessedIndex
		end := lastProcessedIndex + int(poolingIntervalSize)

		result[outIndex] = reducer(slice[start:end]...)
		outIndex++
		lastProcessedIndex = end
	}

	lastBigValue := T(0)

	for i := len(result) - 1; i >= 0; i-- {
		if result[i] == 0 {
			result[i] = lastBigValue
		}

		if result[i] > 0 {
			lastBigValue = result[i]
		}
	}

	return result
}

func GetHammingValues(size int) []float64 {
	hamming := make([]float64, size)

	for i := 0; i < size; i++ {
		hamming[i] = 0.54 - 0.4*math.Cos(2*math.Pi*float64(i)/(float64(size)-1))
	}

	return hamming
}
