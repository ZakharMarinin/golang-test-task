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

func TestUseCase_PutNumber_Success(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	mockStorage.EXPECT().
		PutNumber(mock.Anything, 42).
		Return(nil).
		Once()

	err := useCase.PutNumber(context.Background(), 42)

	assert.NoError(t, err)
}

func TestUseCase_PutNumber_StorageError(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	expectedErr := errors.New("database write failed")
	mockStorage.EXPECT().
		PutNumber(mock.Anything, 42).
		Return(expectedErr).
		Once()

	err := useCase.PutNumber(context.Background(), 42)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestUseCase_PutNumber_PositiveNumber(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	mockStorage.EXPECT().
		PutNumber(mock.Anything, 100).
		Return(nil).
		Once()

	err := useCase.PutNumber(context.Background(), 100)

	assert.NoError(t, err)
}

func TestUseCase_PutNumber_Zero(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	mockStorage.EXPECT().
		PutNumber(mock.Anything, 0).
		Return(nil).
		Once()

	err := useCase.PutNumber(context.Background(), 0)

	assert.NoError(t, err)
}

func TestUseCase_PutNumber_NegativeNumber(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	mockStorage.EXPECT().
		PutNumber(mock.Anything, -42).
		Return(nil).
		Once()

	err := useCase.PutNumber(context.Background(), -42)

	assert.NoError(t, err)
}

func TestUseCase_PutNumber_LargeNumber(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	largeNum := 2147483647
	mockStorage.EXPECT().
		PutNumber(mock.Anything, largeNum).
		Return(nil).
		Once()

	err := useCase.PutNumber(context.Background(), largeNum)

	assert.NoError(t, err)
}

func TestUseCase_PutNumber_ContextPassed(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	var capturedCtx context.Context
	mockStorage.EXPECT().
		PutNumber(mock.Anything, 42).
		Run(func(ctx context.Context, num int) {
			capturedCtx = ctx
		}).
		Return(nil).
		Once()

	ctx := context.WithValue(context.Background(), "request-id", "12345")

	err := useCase.PutNumber(ctx, 42)

	assert.NoError(t, err)
	assert.NotNil(t, capturedCtx)
	assert.Equal(t, "12345", capturedCtx.Value("request-id"))
}

func TestUseCase_PutNumber_TableDriven(t *testing.T) {
	testCases := []struct {
		name   string
		number int
	}{
		{"zero", 0},
		{"positive small", 5},
		{"positive large", 999999},
		{"negative small", -5},
		{"negative large", -999999},
		{"max int32", 2147483647},
		{"min int32", -2147483648},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockStorage := mocks.NewStorage(t)
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

			useCase := &UseCase{
				Storage: mockStorage,
				log:     logger,
			}

			mockStorage.EXPECT().
				PutNumber(mock.Anything, tc.number).
				Return(nil).
				Once()

			err := useCase.PutNumber(context.Background(), tc.number)

			assert.NoError(t, err)
		})
	}
}

func TestUseCase_PutNumber_MultipleCalls(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	mockStorage.EXPECT().
		PutNumber(mock.Anything, 1).
		Return(nil).
		Once()

	mockStorage.EXPECT().
		PutNumber(mock.Anything, 2).
		Return(nil).
		Once()

	mockStorage.EXPECT().
		PutNumber(mock.Anything, 3).
		Return(nil).
		Once()

	assert.NoError(t, useCase.PutNumber(context.Background(), 1))
	assert.NoError(t, useCase.PutNumber(context.Background(), 2))
	assert.NoError(t, useCase.PutNumber(context.Background(), 3))
}

func TestUseCase_PutNumber_LogsError(t *testing.T) {
	// Arrange
	mockStorage := mocks.NewStorage(t)

	// Используем тихий логгер для проверки
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	useCase := &UseCase{
		Storage: mockStorage,
		log:     logger,
	}

	expectedErr := errors.New("storage error")
	mockStorage.EXPECT().
		PutNumber(mock.Anything, 42).
		Return(expectedErr).
		Once()

	err := useCase.PutNumber(context.Background(), 42)

	assert.Error(t, err)
}
