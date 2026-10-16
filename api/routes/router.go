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
	"sl-api/api/routes/list"
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
		listAPI := listHandlers.New(db, validator)

		r.Use(middleware.SetDefaultV1Headers)

		// Public routes
		r.Post("/users", userAPI.Register)
		r.With(middleware.LoginRateLimiter()).Post("/users/login", userAPI.Login)

		r.Group(func(authedRouter chi.Router) {
			authedRouter.Use(middleware.RequireAuth(authSessionStore, authMaxAge))

			authedRouter.Get("/products", productAPI.List)
			authedRouter.Get("/products/{id}", productAPI.Get)
			authedRouter.Post("/products", productAPI.Create)

			authedRouter.Get("/users/{id}", userAPI.Get)
			authedRouter.Post("/users/logout", userAPI.Logout)

			authedRouter.Post("/groups", groupAPI.Create)
			authedRouter.Get("/groups", groupAPI.ListUserGroups)
			authedRouter.Get("/groups/{groupId}/members", groupAPI.ListGroupMembers)
			authedRouter.Get("/invites", inviteAPI.ListUserInvites)
			authedRouter.Post("/groups/{groupId}/invites", inviteAPI.Create)
			authedRouter.Get("/groups/{groupId}/invites", inviteAPI.ListGroupInvites)
			authedRouter.Post("/groups/{groupId}/invites/accept", inviteAPI.AcceptInvite)

			authedRouter.Post("/groups/{groupId}/lists", listAPI.Create)
			authedRouter.Get("/groups/{groupId}/lists", listAPI.List)
			authedRouter.Patch("/lists/{listId}/status", listAPI.UpdateListStatus)
			authedRouter.Post("/lists/{listId}/items", listAPI.AddItem)
			authedRouter.Get("/lists/{listId}/items", listAPI.ListItems)
			authedRouter.Patch("/lists/{listId}/items/{itemId}", listAPI.UpdateItem)
			authedRouter.Delete("/lists/{listId}/items/{itemId}", listAPI.DeleteItem)
			authedRouter.Post("/lists/{listId}/items/{itemId}/increment", listAPI.IncreaseItemQuantity)
			authedRouter.Post("/lists/{listId}/items/{itemId}/decrement", listAPI.DecreaseItemQuantity)
			authedRouter.Patch("/lists/{listId}/items/{itemId}/status", listAPI.UpdateItemStatus)
		})
	})

	return r
}
