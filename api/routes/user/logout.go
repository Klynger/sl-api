package userHandlers

import (
	"encoding/json"
	"net/http"

	"sl-api/api/routes/common/errors"
)

func (a *API) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := a.authSessionStore.Get(r, "auth-session")
	if err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeSessionFailure, "could not read the session"))
		return
	}

	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeSessionSaveFailure, "could not save the session"))
		return
	}

	resp := &LogoutResponse{
		Message: "Logout successful",
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}

type LogoutResponse struct {
	Message string `json:"message"`
}
