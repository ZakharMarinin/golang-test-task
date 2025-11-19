package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"testovoe/internal/domain"
)

//go:generate mockery --name=UseCase --output=mocks/ --outpkg=mocks
type UseCase interface {
	GetSlices(ctx context.Context) ([]int, error)
	PutNumber(ctx context.Context, number int) error
}

type HTTPHandler struct {
	useCase UseCase
	log     *slog.Logger
}

func NewHTTPHandler(useCase UseCase) *HTTPHandler {
	return &HTTPHandler{useCase: useCase}
}

func (h *HTTPHandler) HandleRequest(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.HandleRequest"

		w.Header().Set("Content-Type", "application/json")

		var userNum domain.UserNum

		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.log.Error("Can't read body", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &userNum)
		if err != nil {
			h.log.Error("Can't parse body", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = h.useCase.PutNumber(ctx, userNum.Num)
		if err != nil {
			h.log.Error("could not put num", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		numbers, err := h.useCase.GetSlices(ctx)
		if err != nil {
			h.log.Error("could not get numbers", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, err := json.Marshal(numbers)
		if err != nil {
			h.log.Error("could not marshal response", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = w.Write(response)
		if err != nil {
			h.log.Error("could not write response", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
