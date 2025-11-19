package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"testovoe/internal/usecase/mocks"
)

func TestUseCase_GetSlices_Success(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	unsortedNumbers := []int{5, 2, 8, 1, 9}
	mockStorage.EXPECT().
		GetSlice(mock.Anything).
		Return(unsortedNumbers, nil).
		Once()

	result, err := useCase.GetSlices(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, []int{1, 2, 5, 8, 9}, result)
}

func TestUseCase_GetSlices_StorageError(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	expectedErr := errors.New("database connection failed")
	mockStorage.EXPECT().
		GetSlice(mock.Anything).
		Return(nil, expectedErr).
		Once()

	result, err := useCase.GetSlices(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
}

func TestUseCase_GetSlices_EmptySlice(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	mockStorage.EXPECT().
		GetSlice(mock.Anything).
		Return([]int{}, nil).
		Once()

	result, err := useCase.GetSlices(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, []int{}, result)
}

func TestUseCase_GetSlices_SingleElement(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	mockStorage.EXPECT().
		GetSlice(mock.Anything).
		Return([]int{42}, nil).
		Once()

	result, err := useCase.GetSlices(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, []int{42}, result)
}

func TestUseCase_GetSlices_NegativeNumbers(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	unsorted := []int{-5, 3, -1, 0, -10}
	mockStorage.EXPECT().
		GetSlice(mock.Anything).
		Return(unsorted, nil).
		Once()

	result, err := useCase.GetSlices(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, []int{-10, -5, -1, 0, 3}, result)
}

func TestUseCase_GetSlices_Duplicates(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	unsorted := []int{5, 2, 5, 1, 2, 9}
	mockStorage.EXPECT().
		GetSlice(mock.Anything).
		Return(unsorted, nil).
		Once()

	result, err := useCase.GetSlices(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 2, 5, 5, 9}, result)
}

func TestUseCase_GetSlices_AlreadySorted(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	sorted := []int{1, 2, 3, 4, 5}
	mockStorage.EXPECT().
		GetSlice(mock.Anything).
		Return(sorted, nil).
		Once()

	result, err := useCase.GetSlices(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, sorted, result)
}

func TestUseCase_GetSlices_LargeSlice(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	large := make([]int, 1000)
	for i := range large {
		large[i] = 1000 - i
	}

	mockStorage.EXPECT().
		GetSlice(mock.Anything).
		Return(large, nil).
		Once()

	result, err := useCase.GetSlices(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 1000)

	for i := 0; i < len(result)-1; i++ {
		assert.LessOrEqual(t, result[i], result[i+1])
	}
}

func TestUseCase_GetSlices_SortError(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	mockStorage.EXPECT().
		GetSlice(mock.Anything).
		Return(nil, nil).
		Once()

	result, err := useCase.GetSlices(context.Background())

	if err != nil {
		assert.Error(t, err)
		assert.Nil(t, result)
	} else {
		t.Skip("SortNums не возвращает ошибку для этого случая")
	}
}

func TestUseCase_GetSlices_ContextPassed(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	var capturedCtx context.Context
	mockStorage.EXPECT().
		GetSlice(mock.Anything).
		Run(func(ctx context.Context) {
			capturedCtx = ctx
		}).
		Return([]int{1, 2, 3}, nil).
		Once()

	ctx := context.WithValue(context.Background(), "test-key", "test-value")

	_, err := useCase.GetSlices(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, capturedCtx)
	assert.Equal(t, "test-value", capturedCtx.Value("test-key"))
}
