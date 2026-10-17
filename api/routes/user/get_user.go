package userHandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"sl-api/api/middleware"
	errorsCommon "sl-api/api/routes/common/errors"
)

func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewError(errorsCommon.CodeInvalidUUID, "the provided id is not a valid UUID"))
		return
	}

	authedUserID, err := middleware.GetAuthedUserIDFromCtx(r.Context())
	if err != nil {
		errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
		return
	}

	user, err := a.repository.Read(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorsCommon.NotFound(w, errorsCommon.NewError(errorsCommon.CodeUserNotFound, "user not found"))
			return
		}

		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeDBAccessFailure, "could not read the user"))
		return
	}

	// The full record goes to the user themselves or a group co-member; everyone
	// else gets a trimmed public profile, so the directory cannot be harvested.
	full := id == authedUserID
	if !full {
		shares, err := a.memberRepository.SharesGroup(authedUserID, id)
		if err != nil {
			errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeDBAccessFailure, "could not read the user"))
			return
		}

		full = shares
	}

	var payload any = user.ToPublicDto()
	if full {
		payload = user.ToDto()
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}
