package groupHandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"sl-api/api/middleware"
	"sl-api/api/model/group"
	"sl-api/api/routes/common/errors"
	"sl-api/api/services/group"
	validatorUtil "sl-api/util/validator"
)

// Create godoc
//
//	@summary        Create group
//	@description    Create a new group with the authenticated user as owner
//	@tags           groups
//	@accept         json
//	@produce        json
//	@param          body    body    groupModel.Form    true    "Group form"
//	@success        201     {object}    groupService.CreateOutput
//	@failure        400     {object}    errorsCommon.ErrorResponse
//	@failure        401     {object}    errorsCommon.ErrorResponse
//	@failure        422     {object}    errorsCommon.ErrorResponse
//	@failure        500     {object}    errorsCommon.ErrorResponse
//	@router         /groups [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	form := &groupModel.Form{}
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

	input := groupService.CreateInput{
		Name: form.Name,
	}

	output, err := a.createGroupService.Execute(r.Context(), input)
	if err != nil {
		if errors.Is(err, middleware.ErrUnauthorized) {
			errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
			return
		}

		errorsCommon.WriteDBError(w, err, errorsCommon.NewError(errorsCommon.CodeDBInsertFailure, "could not create the group"))
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(output); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}
