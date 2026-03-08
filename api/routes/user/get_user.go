package userHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"sl-api/api/routes/common/errors"
)

func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.RespInvalidURLParamID)
		return
	}

	user, err := a.repository.Read(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	dto := user.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		errorsCommon.ServerError(w, errorsCommon.RespJSONEncodeFailure)
		return
	}
}
