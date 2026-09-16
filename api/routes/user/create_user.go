package userHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"sl-api/api/model/user"
	"sl-api/api/routes/common/errors"
	validatorUtil "sl-api/util/validator"
)

// Register godoc
//
//	@summary		Register user
//	@description	Register user
//	@tags			users
//	@accept			json
//	@product		json
//	@param			body	body	userModel.Form	true	"User form"
//	@success		201
//	@failure		400	{object}	errorsCommon.ErrorResponse
//	@failure		409	{object}	errorsCommon.ErrorResponse
//	@failure		422	{object}	errorsCommon.ErrorResponse
//	@failure		500	{object}	errorsCommon.ErrorResponse
//	@router			/users [post]
func (a *API) Register(w http.ResponseWriter, r *http.Request) {
	form := &userModel.Form{}
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

	newUser, err := form.ToModel()
	if err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeInternalError, "could not process the password"))
		return
	}
	newUser.ID = uuid.New()

	_, err = a.repository.Create(newUser)
	if err != nil {
		if item, ok := errorsCommon.ClassifyDBError(err); ok && item.Code == errorsCommon.CodeDuplicateEntry {
			errorsCommon.Conflict(w, errorsCommon.NewError(errorsCommon.CodeUsernameTaken, "username is already taken"))
			return
		}

		errorsCommon.WriteDBError(w, err, errorsCommon.NewError(errorsCommon.CodeDBInsertFailure, "could not create the user"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
