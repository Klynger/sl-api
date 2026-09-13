package inviteHandlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"sl-api/api/middleware"
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
//	@param          body    body    groupService.CreateInviteInput    true    "Create invite form"
//	@success        201     {object}    groupService.CreateInviteOutput
//	@failure        400     {object}    errorsCommon.ErrorResponse
//	@failure        401     {object}    errorsCommon.ErrorResponse
//	@failure        403     {object}    errorsCommon.ErrorResponse
//	@failure        404     {object}    errorsCommon.ErrorResponse
//	@failure        409     {object}    errorsCommon.ErrorResponse
//	@failure        422     {object}    errorsCommon.ErrorResponse
//	@failure        500     {object}    errorsCommon.ErrorResponse
//	@router         /invites [post]
func (api *API) Create(w http.ResponseWriter, r *http.Request) {
	form := &inviteModel.Form{}

	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewError(errorsCommon.CodeInvalidJSON, "invalid JSON in request body"))
		return
	}

	if err := api.validator.Struct(form); err != nil {
		items := validatorUtil.ToErrResponse(err)
		if items == nil {
			errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeInternalError, "unexpected validation error"))
			return
		}

		errorsCommon.ValidationErrors(w, items...)
		return
	}

	groupID, err := uuid.Parse(chi.URLParam(r, "groupId"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewErrorWithMeta(errorsCommon.CodeInvalidUUID, "the provided group id is not a valid UUID", map[string]any{"field": "groupId"}))
		return
	}

	invitedUserID, err := uuid.Parse(form.InvitedUserID)
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewErrorWithMeta(errorsCommon.CodeInvalidUUID, "the provided invited user id is not a valid UUID", map[string]any{"field": "invitedUserId"}))
		return
	}

	input := groupService.CreateInviteInput{
		GroupID:       groupID,
		InvitedUserID: invitedUserID,
	}

	output, err := api.inviteService.Execute(r.Context(), input)

	if err != nil {
		switch {
		case errors.Is(err, middleware.ErrUnauthorized):
			errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
		case errors.Is(err, groupService.ErrInviteSelf):
			errorsCommon.BadRequest(w, errorsCommon.NewError(errorsCommon.CodeInviteSelf, "cannot invite yourself"))
		case errors.Is(err, groupService.ErrGroupNotFound):
			errorsCommon.NotFound(w, errorsCommon.NewError(errorsCommon.CodeGroupNotFound, "group not found"))
		case errors.Is(err, groupService.ErrAlreadyInvited):
			errorsCommon.Conflict(w, errorsCommon.NewError(errorsCommon.CodeAlreadyInvited, "user already invited"))
		case errors.Is(err, groupService.ErrAlreadyMember):
			errorsCommon.Conflict(w, errorsCommon.NewError(errorsCommon.CodeAlreadyMember, "user is already a member"))
		case errors.Is(err, groupService.ErrNotAMember):
			errorsCommon.Forbidden(w, errorsCommon.NewError(errorsCommon.CodeNotAMember, "you are not a member of this group"))
		case errors.Is(err, groupService.ErrInsufficientPermissions):
			errorsCommon.Forbidden(w, errorsCommon.NewError(errorsCommon.CodeInsufficientPermissions, "only group owners can invite"))
		default:
			fmt.Println("Error executing invite service:", err)
			errorsCommon.WriteDBError(w, err, errorsCommon.NewError(errorsCommon.CodeDBInsertFailure, "could not create the invite"))
		}
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(output); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}
