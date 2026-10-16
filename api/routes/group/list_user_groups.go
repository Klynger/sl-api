package groupHandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"sl-api/api/middleware"
	errorsCommon "sl-api/api/routes/common/errors"
)

func (a *API) ListUserGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := a.groupService.ListUserGroups(r.Context())
	if err != nil {
		switch {
		case errors.Is(err, middleware.ErrUnauthorized):
			errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
		default:
			errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeDBAccessFailure, "could not list the groups"))
		}
		return
	}

	if err := json.NewEncoder(w).Encode(groups.ToDto()); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}
