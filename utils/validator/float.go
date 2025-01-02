package validator

func FloatIsPositive(i float64) bool {
	return i > 0
}

func FloatIsNegative(i float64) bool {
	return i < 0
}

func FloatIsZero(i float64) bool {
	return i == 0.0
}

func FloatIsGreaterThan(i float64, j float64) bool {
	return i > j
}

func FloatIsLessThan(i float64, j float64) bool {
	return i < j
}

func FloatBetween(i, min, max float64) bool {
	return i >= min && i <= max
}
