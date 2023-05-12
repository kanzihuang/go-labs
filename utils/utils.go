package utils

import "math"

func Min(numbers ...int) int {
	result := math.MaxInt
	for _, num := range numbers {
		if num < result {
			result = num
		}
	}
	return result
}

func Max(numbers ...int) int {
	result := math.MinInt
	for _, num := range numbers {
		if num > result {
			result = num
		}
	}
	return result
}
