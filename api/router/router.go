package router

import (
	"github.com/go-chi/chi/v5"

	"sl-api/api/resource/health"
	"sl-api/api/resource/product"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/livez", health.Read)

	r.Route("/v1", func(v1Router chi.Router) {
		productAPI := &product.API{}

		v1Router.Get("/products/{id}", productAPI.Read)
		v1Router.Post("/products", productAPI.Create)
	})

	return r
}
