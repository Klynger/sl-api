package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"

	"sl-api/api/middleware"
	"sl-api/api/resource/health"
	"sl-api/api/resource/product"
	"sl-api/api/resource/user"
)

func New(db *gorm.DB, validator *validator.Validate, authSessionStore *sessions.CookieStore, isDebugging bool, authMaxAge int) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/livez", health.Read)

	r.Route("/v1", func(r chi.Router) {
		productAPI := product.New(db, validator)
		userAPI := user.New(db, authSessionStore, authMaxAge, isDebugging, validator)

		// Public routes
		r.Post("/users", userAPI.Register)
		r.Post("/users/login", userAPI.Login)

		r.Group(func(authedRouter chi.Router) {
			authedRouter.Use(middleware.RequireAuth(authSessionStore, authMaxAge))

			authedRouter.Get("/products", productAPI.List)
			authedRouter.Get("/products/{id}", productAPI.Read)
			authedRouter.Post("/products", productAPI.Create)

			authedRouter.Get("/users/{id}", userAPI.Read)
			authedRouter.Post("/users/logout", userAPI.Logout)
		})
	})

	return r
}
