package groupHandlers

import (
	"encoding/json"
	"net/http"

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
//	@success        201     {object}    groupModel.DTO
//	@failure        400     {object}    err.Error
//	@failure        401     {object}    err.Error
//	@failure        422     {object}    err.Errors
//	@failure        500     {object}    err.Error
//	@router         /groups [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	form := &groupModel.Form{}
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

	input := groupService.CreateInput{
		Name: form.Name,
	}

	output, err := a.createGroupService.Execute(r.Context(), input)
	if err != nil {
		errorsCommon.ServerError(w, errorsCommon.RespDBDataInsertFailure)
		return
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		errorsCommon.ServerError(w, errorsCommon.RespJSONEncodeFailure)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
