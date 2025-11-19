package usecase

import (
	"context"
)

func (u *UseCase) GetSlices(ctx context.Context) ([]int, error) {
	const op = "useCase.GetSlices"

	numbers, err := u.Storage.GetSlice(ctx)
	if err != nil {
		u.log.Error("failed to get slices", op, err)
		return nil, err
	}

	numbers, err = SortNums(numbers)
	if err != nil {
		u.log.Error("failed to sort numbers", op, err)
		return nil, err
	}

	return numbers, nil
}
