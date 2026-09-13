package userHandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"sl-api/api/routes/common/errors"
)

func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewError(errorsCommon.CodeInvalidUUID, "the provided id is not a valid UUID"))
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

	dto := user.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}
