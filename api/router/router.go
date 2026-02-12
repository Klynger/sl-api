package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"

	"sl-api/api/resource/health"
	"sl-api/api/resource/product"
	"sl-api/api/resource/user"
)

func New(db *gorm.DB, authSessionStore *sessions.CookieStore, v *validator.Validate) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/livez", health.Read)

	r.Route("/v1", func(r chi.Router) {
		productAPI := product.New(db, v)
		userAPI := user.New(db, authSessionStore, v)

		r.Get("/products", productAPI.List)
		r.Get("/products/{id}", productAPI.Read)
		r.Post("/products", productAPI.Create)

		r.Post("/users", userAPI.Register)
		r.Get("/users/{id}", userAPI.Read)
		r.Post("/users/login", userAPI.Login)
	})

	return r
}
