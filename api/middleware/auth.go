package middleware

import (
	"github.com/gorilla/sessions"
	"net/http"

	e "sl-api/api/resource/common/err"
)

func RequireAuth(authSessionStore *sessions.CookieStore, maxAge int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := authSessionStore.Get(r, "auth-session")
			if err != nil {
				e.ServerError(w, e.RespSessionAccessFailure)
			}

			userID, ok := session.Values["user_id"].(string)
			if !ok || userID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Reset session expiration on each request
			session.Options.MaxAge = maxAge
			err = session.Save(r, w)
			if err != nil {
				e.ServerError(w, e.RespGenericFailure)
				return
			}

			// User is authenticated, proceed
			next.ServeHTTP(w, r)
		})
	}
}
