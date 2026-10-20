package v1

import (
	"net/http"
	"net/netip"

	docv1 "github.com/Weit145/template-go-monolit/docs/api/v1"
	"github.com/Weit145/template-go-monolit/internal/transport/api/v1/handler"
	"github.com/Weit145/template-go-monolit/internal/transport/api/v1/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type UseCase interface {
	handler.HealthService
	handler.UserService
	handler.WebSocketService
}

func NewRouter(ser UseCase, trustedProxies []netip.Prefix) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recovery, middleware.RequestID, middleware.Logging)
	r.Get("/ws", handler.WebSocket(ser))
	r.Get("/_health", handler.Health(ser))
	r.Post("/api/v1/users/create-and-delete", handler.CreateAndDeleteUser(ser))
	r.Get("/swagger/doc.yaml", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(docv1.Spec)
	})
	r.Handle("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.yaml")))

	return r
}
