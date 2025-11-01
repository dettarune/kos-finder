// routes/routes.go
package routes

import (
	"net/http"

	"github.com/dettarune/kos-finder/internal/delivery/handler"
	"github.com/gorilla/mux"
)

type RouteConfig struct {
	Router         *mux.Router
	UserHandler    *handler.UserHandler
	ProductHandler *handler.ProductHandler
}

func NewRouterConfig(userHandler *handler.UserHandler, productHandler *handler.ProductHandler) *RouteConfig {
	return &RouteConfig{
		Router:         mux.NewRouter(),
		UserHandler:    userHandler,
		ProductHandler: productHandler,
	}
}

func (c *RouteConfig) SetupRoutes() {
	c.Router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello 2205"))
	}).Methods("GET")

	c.Router.HandleFunc("/api/register", c.UserHandler.RegisterHandler).Methods("POST")
	c.Router.HandleFunc("/api/login", c.UserHandler.LoginHandler).Methods("POST")

	// c.Router.HandleFunc("/api/products", c.ProductHandler.CreateProductHandler).Methods("POST")
	// c.Router.HandleFunc("/api/products", c.ProductHandler.GetProductsHandler).Methods("GET")
}
