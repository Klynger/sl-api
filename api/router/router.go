package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"sl-api/api/resource/health"
	"sl-api/api/resource/product"
)

func New(db *gorm.DB, v *validator.Validate) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/livez", health.Read)

	r.Route("/v1", func(r chi.Router) {
		productAPI := product.New(db, v)

		r.Get("/products", productAPI.List)
		r.Get("/products/{id}", productAPI.Read)
		r.Post("/products", productAPI.Create)
	})

	return r
}
