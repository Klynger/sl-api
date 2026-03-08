package userHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
	"net/http"

	"sl-api/api/model/user"
	"sl-api/api/routes/common/errors"
	validatorUtil "sl-api/util/validator"
)

func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	form := &userModel.LoginForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		errorsCommon.ServerError(w, errorsCommon.RespJSONDecodeFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			errorsCommon.ServerError(w, errorsCommon.RespJSONEncodeFailure)
			return
		}

		errorsCommon.ValidationErrors(w, respBody)
		return
	}

	userCredentials, err := a.repository.ReadByUsername(form.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// TODO: Do something to avoid timing attacks
			w.WriteHeader(http.StatusNotFound)
			return
		}

		errorsCommon.ServerError(w, errorsCommon.RespDBDataAccessFailure)
		return
	}

	validationResult := userCredentials.ValidateCredentials(form)

	if !validationResult {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session, err := a.authSessionStore.Get(r, "auth-session")
	if err != nil {
		errorsCommon.ServerError(w, errorsCommon.RespSessionAccessFailure)
		return
	}

	session.Values["user_id"] = userCredentials.ID.String()
	session.Values["username"] = userCredentials.Username

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   a.authMaxAge,
		HttpOnly: !a.isDebugging, // Allow JS access to cookies in debug mode for testing purposes
		Secure:   !a.isDebugging, // Don't require secure cookies in debug mode for testing purposes
	}

	err = session.Save(r, w)
	if err != nil {
		fmt.Println("Session save error: ", err)
		errorsCommon.ServerError(w, errorsCommon.RespGenericFailure)
		return
	}

	resp := &LoginResponse{
		Message:  "Login successful",
		Username: userCredentials.Username,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		errorsCommon.ServerError(w, errorsCommon.RespJSONEncodeFailure)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type LoginResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}
