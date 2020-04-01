package mathutil

// MinInt returns the max of the given values.
func MinInt(a, b int) int {
	if b < a {
		return b
	}
	return a
}

// MinInt32 returns the max of the given values.
func MinInt32(a, b int32) int32 {
	if b < a {
		return b
	}
	return a
}

// MinInt64 returns the max of the given values.
func MinInt64(a, b int64) int64 {
	if b < a {
		return b
	}
	return a
}

// MinFloat32 returns the max of the given values.
func MinFloat32(a, b float32) float32 {
	if b < a {
		return b
	}
	return a
}

// MinFloat64 returns the max of the given values.
func MinFloat64(a, b float64) float64 {
	if b < a {
		return b
	}
	return a
}

// MaxInt returns the max of the given values.
func MaxInt(a, b int) int {
	if b > a {
		return b
	}
	return a
}

// MaxInt32 returns the max of the given values.
func MaxInt32(a, b int32) int32 {
	if b > a {
		return b
	}
	return a
}

// MaxInt64 returns the max of the given values.
func MaxInt64(a, b int64) int64 {
	if b > a {
		return b
	}
	return a
}

// MaxFloat32 returns the max of the given values.
func MaxFloat32(a, b float32) float32 {
	if b > a {
		return b
	}
	return a
}

// MaxFloat64 returns the max of the given values.
func MaxFloat64(a, b float64) float64 {
	if b > a {
		return b
	}
	return a
}

// ClampInt clamps the given value to the given range.
func ClampInt(x, lower, upper int) int {
	return MinInt(upper, MaxInt(lower, x))
}

// ClampInt32 clamps the given value to the given range.
func ClampInt32(x, lower, upper int32) int32 {
	return MinInt32(upper, MaxInt32(lower, x))
}

// ClampInt64 clamps the given value to the given range.
func ClampInt64(x, lower, upper int64) int64 {
	return MinInt64(upper, MaxInt64(lower, x))
}

// ClampFloat32 clamps the given value to the given range.
func ClampFloat32(x, lower, upper float32) float32 {
	return MinFloat32(upper, MaxFloat32(lower, x))
}

// ClampFloat64 clamps the given value to the given range.
func ClampFloat64(x, lower, upper float64) float64 {
	return MinFloat64(upper, MaxFloat64(lower, x))
}
