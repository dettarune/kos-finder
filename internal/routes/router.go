// routes/routes.go
package routes

import (
	"net/http"

	"github.com/dettarune/kos-finder/internal/delivery/handler"
	"github.com/dettarune/kos-finder/internal/middleware"
	"github.com/gorilla/mux"
)

type RouteConfig struct {
	Router         *mux.Router
	UserHandler    *handler.UserHandler
	KosHandler *handler.KosHandler
	AuthMiddleware *middleware.AuthMiddleware
}

func NewRouterConfig(userHandler *handler.UserHandler, productHandler *handler.KosHandler, authmiddleware *middleware.AuthMiddleware) *RouteConfig {
	return &RouteConfig{
		Router:         mux.NewRouter(),
		UserHandler:    userHandler,
		KosHandler: productHandler,
		AuthMiddleware: authmiddleware,
	}
}

func (r *RouteConfig) SetupGuestRoutes() {
	r.Router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello 2205"))
	}).Methods("GET")

	r.Router.HandleFunc("/api/auth/register", r.UserHandler.RegisterHandler).Methods("POST")
	r.Router.HandleFunc("/api/auth/login", r.UserHandler.LoginHandler).Methods("POST")
	r.Router.HandleFunc("/api/auth/verify", r.UserHandler.VerifyHandler).Methods("GET")
}

func (r *RouteConfig) SetupAuthRoutes() {
	protected := r.Router.PathPrefix("/api").Subrouter()

	protected.Use(r.AuthMiddleware.Authenticate)

	// protected.HandleFunc("/kos",r.KosHandler. )

}