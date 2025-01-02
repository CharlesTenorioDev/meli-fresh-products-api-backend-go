package validator

func IntIsPositive(i int) bool {
	return i > 0
}

func IntIsNegative(i int) bool {
	return i < 0
}

func IntIsZero(i int) bool {
	return i == 0
}

func IntIsGreaterThan(i int, j int) bool {
	return i > j
}

func IntIsLessThan(i int, j int) bool {
	return i < j
}

func IntBetween(i, min, max int) bool {
	return i >= min && i <= max
}
