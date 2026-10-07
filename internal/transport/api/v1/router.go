package v1

import (
	"net/http"
	"net/netip"

	"github.com/Weit145/template-go-monolit/internal/transport/api/v1/handler"
	"github.com/Weit145/template-go-monolit/internal/transport/api/v1/middleware"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type UseCase interface {
	handler.HealthService
}

func NewRouter(ser UseCase, trustedProxies []netip.Prefix) *mux.Router {

	r := mux.NewRouter()
	r.Use(middleware.Recovery, middleware.RequestID, middleware.Logging)
	r.HandleFunc("/_health", handler.Health(ser)).Methods(http.MethodGet)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}
