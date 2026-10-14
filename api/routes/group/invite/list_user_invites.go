package inviteHandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"sl-api/api/middleware"
	"sl-api/api/routes/common/errors"
)

func (api *API) ListUserInvites(w http.ResponseWriter, r *http.Request) {
	invites, err := api.listUserInvitesService.Execute(r.Context())
	if err != nil {
		switch {
		case errors.Is(err, middleware.ErrUnauthorized):
			errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
		default:
			errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeDBAccessFailure, "could not list the invites"))
		}
		return
	}

	if err := json.NewEncoder(w).Encode(invites.ToDto()); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}
