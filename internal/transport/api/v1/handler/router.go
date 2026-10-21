package handlers

import (
	"net/http"

	docv1 "github.com/Weit145/template-go-monolit/docs/api/v1"
	"github.com/Weit145/template-go-monolit/internal/transport/api/v1/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler interface {
	Routes() http.Handler
}

type UseCase interface {
	HealthService
	UserService
	WebSocketService
}

type handler struct {
	ser UseCase
}

func NewHandler(
	ser UseCase,
) Handler {
	return &handler{
		ser: ser,
	}
}

func (h *handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recovery, middleware.RequestID, middleware.Logging)
	r.Get("/ws", h.WebSocket)
	r.Get("/_health", h.Health)
	r.Post("/api/v1/users/create-and-delete", h.CreateAndDeleteUser)
	r.Get("/swagger/doc.yaml", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(docv1.Spec)
	})
	r.Handle("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.yaml")))

	return r
}
