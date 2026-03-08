package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"

	"sl-api/api/middleware"
	"sl-api/api/routes/group"
	"sl-api/api/routes/group/invite"
	"sl-api/api/routes/health"
	"sl-api/api/routes/product"
	"sl-api/api/routes/user"
)

func New(db *gorm.DB, validator *validator.Validate, authSessionStore *sessions.CookieStore, isDebugging bool, authMaxAge int) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/livez", health.Read)

	r.Route("/v1", func(r chi.Router) {
		productAPI := productHandlers.New(db, validator)
		userAPI := userHandlers.New(db, authSessionStore, authMaxAge, isDebugging, validator)
		groupAPI := groupHandlers.New(db, validator)
		inviteAPI := inviteHandlers.New(db, validator)

		r.Use(middleware.SetDefaultV1Headers)

		// Public routes
		r.Post("/users", userAPI.Register)
		r.Post("/users/login", userAPI.Login)

		r.Group(func(authedRouter chi.Router) {
			authedRouter.Use(middleware.RequireAuth(authSessionStore, authMaxAge))

			authedRouter.Get("/products", productAPI.List)
			authedRouter.Get("/products/{id}", productAPI.Get)
			authedRouter.Post("/products", productAPI.Create)

			authedRouter.Get("/users/{id}", userAPI.Get)
			authedRouter.Post("/users/logout", userAPI.Logout)

			authedRouter.Post("/groups", groupAPI.Create)
			authedRouter.Post("/groups/{groupId}/invites", inviteAPI.Create)
			authedRouter.Post("/groups/{groupId}/invites/accept", inviteAPI.AcceptInvite)
		})
	})

	return r
}
