package userHandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"gorm.io/gorm"

	"sl-api/api/middleware"
	errorsCommon "sl-api/api/routes/common/errors"
)

// Me godoc
//
//	@summary		Get the authenticated user
//	@description	Return the full record of the user the session belongs to
//	@tags			users
//	@produce		json
//	@success		200	{object}	userModel.DTO
//	@failure		401	{object}	errorsCommon.ErrorResponse
//	@failure		500	{object}	errorsCommon.ErrorResponse
//	@router			/users/me [get]
func (a *API) Me(w http.ResponseWriter, r *http.Request) {
	authedUserID, err := middleware.GetAuthedUserIDFromCtx(r.Context())
	if err != nil {
		errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
		return
	}

	user, err := a.repository.Read(authedUserID)
	if err != nil {
		// A session can outlive its user (the user was deleted after logging in).
		// That session no longer identifies anyone, so it is treated as signed out.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
			return
		}

		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeDBAccessFailure, "could not read the user"))
		return
	}

	if err := json.NewEncoder(w).Encode(user.ToDto()); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}
