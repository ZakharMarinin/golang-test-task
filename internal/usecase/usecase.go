package usecase

import (
	"context"
	"log/slog"
)

//go:generate mockery --name=UseCase --output=mocks/ --outpkg=mocks
type Storage interface {
	PutNumber(ctx context.Context, num int) error
	GetSlice(ctx context.Context) (numbers []int, err error)
}

type UseCase struct {
	log     *slog.Logger
	Storage Storage
}

func NewUseCase(log *slog.Logger, storage Storage) *UseCase {
	return &UseCase{
		log:     log,
		Storage: storage,
	}
}
