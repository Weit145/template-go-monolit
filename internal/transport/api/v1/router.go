package v1

import (
	"net/http"
	"net/netip"

	docv1 "github.com/Weit145/template-go-monolit/docs/api/v1"
	"github.com/Weit145/template-go-monolit/internal/transport/api/v1/handler"
	"github.com/Weit145/template-go-monolit/internal/transport/api/v1/middleware"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type UseCase interface {
	handler.HealthService
	handler.UserService
}

func NewRouter(ser UseCase, trustedProxies []netip.Prefix) *mux.Router {

	r := mux.NewRouter()
	r.Use(middleware.Recovery, middleware.RequestID, middleware.Logging)
	r.HandleFunc("/_health", handler.Health(ser)).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users", handler.CreateUser(ser)).Methods(http.MethodPost)
	r.HandleFunc("/swagger/doc.yaml", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(docv1.Spec)
	}).Methods(http.MethodGet)
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL("/swagger/doc.yaml")))

	return r
}
