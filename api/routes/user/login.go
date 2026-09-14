package userHandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"gorm.io/gorm"

	"sl-api/api/model/user"
	"sl-api/api/routes/common/errors"
	validatorUtil "sl-api/util/validator"
)

func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	form := &userModel.LoginForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewError(errorsCommon.CodeInvalidJSON, "invalid JSON in request body"))
		return
	}

	if err := a.validator.Struct(form); err != nil {
		items := validatorUtil.ToErrResponse(err)
		if items == nil {
			errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeInternalError, "unexpected validation error"))
			return
		}

		errorsCommon.ValidationErrors(w, items...)
		return
	}

	userCredentials, err := a.repository.ReadByUsername(form.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// TODO: Do something to avoid timing attacks
			errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeInvalidCredentials, "invalid username or password"))
			return
		}

		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeDBAccessFailure, "could not read the user"))
		return
	}

	validationResult := userCredentials.ValidateCredentials(form)

	if !validationResult {
		errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeInvalidCredentials, "invalid username or password"))
		return
	}

	session, err := a.authSessionStore.Get(r, "auth-session")
	if err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeSessionFailure, "could not read the session"))
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
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeSessionSaveFailure, "could not save the session"))
		return
	}

	resp := &LoginResponse{
		Message:  "Login successful",
		Username: userCredentials.Username,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}

type LoginResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}
