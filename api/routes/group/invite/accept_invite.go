package inviteHandlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"sl-api/api/middleware"
	"sl-api/api/routes/common/errors"
	"sl-api/api/services/group"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *API) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	groupID, err := uuid.Parse(chi.URLParam(r, "groupId"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewErrorWithMeta(errorsCommon.CodeInvalidUUID, "the provided group id is not a valid UUID", map[string]any{"field": "groupId"}))
		return
	}

	input := groupService.AcceptInviteInput{
		GroupID: groupID,
	}

	output, err := api.acceptInviteService.Execute(r.Context(), input)

	if err != nil {
		switch {
		case errors.Is(err, middleware.ErrUnauthorized):
			errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
		case errors.Is(err, groupService.ErrInviteNotFound):
			errorsCommon.NotFound(w, errorsCommon.NewError(errorsCommon.CodeInviteNotFound, "invite not found"))
		default:
			fmt.Println("Error executing accept invite service:", err)
			errorsCommon.WriteDBError(w, err, errorsCommon.NewError(errorsCommon.CodeDBInsertFailure, "could not accept the invite"))
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)

	if err := json.NewEncoder(w).Encode(output); err != nil {
		fmt.Println("Error encoding accept invite output:", err)
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}
