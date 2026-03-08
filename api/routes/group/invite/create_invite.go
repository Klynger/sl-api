package inviteHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sl-api/api/model/invite"
	"sl-api/api/routes/common/errors"
	"sl-api/api/services/group"
	validatorUtil "sl-api/util/validator"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Create godoc
//
//	@summary        Create invite
//	@description    Create a new invite for a user to join a group
//	@tags           invites
//	@accept         json
//	@produce        json
//	@param          body    body    group_service.CreateInviteInput    true    "Create invite form"
//	@success        201     {object}    group_service.CreateInviteOutput
//	@failure        400     {object}    e.Error
//	@failure        401     {object}    e.Error
//	@failure        422     {object}    e.Errors
//	@failure        500     {object}    e.Error
//	@router         /invites [post]
func (api *API) Create(w http.ResponseWriter, r *http.Request) {
	form := &inviteModel.Form{}

	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		errorsCommon.ServerError(w, errorsCommon.RespJSONDecodeFailure)
		return
	}

	if err := api.validator.Struct(form); err != nil {
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			errorsCommon.ServerError(w, errorsCommon.RespJSONEncodeFailure)
			return
		}

		errorsCommon.ValidationErrors(w, respBody)
		return
	}

	groupID, err := uuid.Parse(chi.URLParam(r, "groupId"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.RespInvalidUUID)
		return
	}

	invitedUserID, err := uuid.Parse(form.InvitedUserID)
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.RespInvalidUUID)
		return
	}

	input := groupService.CreateInviteInput{
		GroupID:       groupID,
		InvitedUserID: invitedUserID,
	}

	output, err := api.inviteService.Execute(r.Context(), input)

	if err != nil {
		fmt.Println("Error executing invite service:", err)

		errorsCommon.ServerError(w, errorsCommon.RespGenericFailure)
		return

		// TODO: Create these errors
		// switch err {
		// case group_service.ErrGroupNotFound:
		// 	e.NotFound(w, e.RespGroupNotFound)
		// case group_service.ErrUserNotFound:
		// 	e.NotFound(w, e.RespUserNotFound)
		// case group_service.ErrAlreadyInvited:
		// 	e.BadRequest(w, e.RespAlreadyInvited)
		// case group_service.ErrAlreadyMember:
		// 	e.BadRequest(w, e.RespAlreadyMember)
		// case group_service.ErrNotGroupOwner:
		// 	e.Unauthorized(w, e.RespNotGroupOwner)
		// default:
		// 	e.ServerError(w, e.RespDBDataInsertFailure)
		// }
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		errorsCommon.ServerError(w, errorsCommon.RespJSONEncodeFailure)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
