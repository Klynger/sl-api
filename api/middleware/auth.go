package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"

	errorsCommon "sl-api/api/routes/common/errors"
)

// ErrUnauthorized is returned when the request context carries no valid
// authenticated user. Services propagate it so handlers can map it to a 401.
var ErrUnauthorized = errors.New("unauthorized")

func RequireAuth(authSessionStore *sessions.CookieStore, maxAge int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := authSessionStore.Get(r, "auth-session")
			if err != nil {
				errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeSessionFailure, "could not read the session"))
				return
			}

			userID, ok := session.Values["user_id"].(string)
			if !ok || userID == "" {
				errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", userID)
			r = r.WithContext(ctx)

			// Reset session expiration on each request
			session.Options.MaxAge = maxAge
			err = session.Save(r, w)
			if err != nil {
				errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeSessionSaveFailure, "could not save the session"))
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
		return uuid.Nil, ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)

	if err != nil {
		return uuid.Nil, ErrUnauthorized
	}

	return userID, nil
}
