package inviteHandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"sl-api/api/middleware"
	"sl-api/api/routes/common/errors"
	"sl-api/api/services/group"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *API) ListGroupInvites(w http.ResponseWriter, r *http.Request) {
	groupID, err := uuid.Parse(chi.URLParam(r, "groupId"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewErrorWithMeta(errorsCommon.CodeInvalidUUID, "the provided group id is not a valid UUID", map[string]any{"field": "groupId"}))
		return
	}

	invites, err := api.groupService.ListGroupInvites(r.Context(), groupID)
	if err != nil {
		switch {
		case errors.Is(err, middleware.ErrUnauthorized):
			errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
		case errors.Is(err, groupService.ErrGroupNotFound):
			// Reported as 404 so the endpoint does not reveal a group the caller cannot see.
			errorsCommon.NotFound(w, errorsCommon.NewError(errorsCommon.CodeGroupNotFound, "group not found"))
		default:
			errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeDBAccessFailure, "could not list the group's invites"))
		}
		return
	}

	if err := json.NewEncoder(w).Encode(invites.ToDto()); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}
