package intutil

import "math"

// Integer limit values.
const (
    MaxInt   = int(^uint(0) >> 1)
    MinInt   = int(-MaxInt - 1)
    MaxInt32 = int32(math.MaxInt32)
    MinInt32 = int32(math.MinInt32)
    MaxInt64 = int64(math.MaxInt64)
    MinInt64 = int64(math.MinInt64)
)

func Max(left, right int) int {
	if left > right {
		return left
	} else {
		return right
	}
}

func Min(left, right int) int {
	if left < right {
		return left
	} else {
		return right
	}
}

// Abs returns the absolute value of x.
func Abs(x int) int {
    switch {
    case x >= 0:
        return x
    case x > MinInt:
        return -x
    }
    panic("math/int.Abs: invalid argument")

}

// Abs32 returns the absolute value of x.
func Abs32(x int32) int32 {
    switch {
    case x >= 0:
        return x
    case x > MinInt32:
        return -x
    }
    panic("math/int.Abs32: invalid argument")

}

// Abs64 returns the absolute value of x.
func Abs64(x int64) int64 {
    switch {
    case x >= 0:
        return x
    case x > MinInt64:
        return -x
    }
    panic("math/int.Abs64: invalid argument")
}