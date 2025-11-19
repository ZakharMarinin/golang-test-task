package usecase

import "context"

func (u *UseCase) PutNumber(ctx context.Context, number int) error {
	const op = "useCase.PutNumber"

	err := u.Storage.PutNumber(ctx, number)
	if err != nil {
		u.log.Error("failed to put number", "op", op, "error", err)
		return err
	}

	return nil
}
