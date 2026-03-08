package middleware

import (
	"context"
	"fmt"
	"net/http"
	errorsCommon "sl-api/api/routes/common/errors"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

func RequireAuth(authSessionStore *sessions.CookieStore, maxAge int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := authSessionStore.Get(r, "auth-session")
			if err != nil {
				errorsCommon.ServerError(w, errorsCommon.RespSessionAccessFailure)
				return
			}

			userID, ok := session.Values["user_id"].(string)
			if !ok || userID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", userID)
			r = r.WithContext(ctx)

			// Reset session expiration on each request
			session.Options.MaxAge = maxAge
			err = session.Save(r, w)
			if err != nil {
				errorsCommon.ServerError(w, errorsCommon.RespGenericFailure)
				return
			}

			// User is authenticated, proceed
			next.ServeHTTP(w, r)
		})
	}
}

func GetAuthedUserIDFromCtx(ctx context.Context) (uuid.UUID, error) {
	userIDStr, ok := ctx.Value("user_id").(string)

	if !ok {
		return uuid.Nil, fmt.Errorf("UNAUTHORIZED")
	}

	userID, err := uuid.Parse(userIDStr)

	if err != nil {
		return uuid.Nil, fmt.Errorf("UNAUTHORIZED")
	}

	return userID, nil
}
