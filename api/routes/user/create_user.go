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
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/users [post]
func (a *API) Register(w http.ResponseWriter, r *http.Request) {
	form := &userModel.Form{}
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

	newUser := form.ToModel()
	newUser.ID = uuid.New()

	_, err := a.repository.Create(newUser)
	if err != nil {
		errorsCommon.ServerError(w, errorsCommon.RespDBDataInsertFailure)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
