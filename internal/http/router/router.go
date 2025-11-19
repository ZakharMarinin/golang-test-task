package router

import (
	"context"
	"testovoe/internal/http/handlers"

	"github.com/go-chi/chi/v5"
)

func Router(ctx context.Context, router *chi.Mux, http *handlers.HTTPHandler) {
	router.Post("/put-num", http.HandleRequest(ctx))
}
