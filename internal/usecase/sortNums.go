package usecase

import (
	"slices"
)

func SortNums(numbers []int) ([]int, error) {
	numbersLen := len(numbers)
	if numbersLen < 2 {
		return numbers, nil
	}

	slices.Sort(numbers)

	return numbers, nil
}
