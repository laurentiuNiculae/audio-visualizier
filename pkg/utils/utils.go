package utils

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
