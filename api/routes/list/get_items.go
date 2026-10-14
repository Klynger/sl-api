package listHandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"sl-api/api/middleware"
	errorsCommon "sl-api/api/routes/common/errors"
	listService "sl-api/api/services/list"
)

func (a *API) ListItems(w http.ResponseWriter, r *http.Request) {
	listID, err := uuid.Parse(chi.URLParam(r, "listId"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewErrorWithMeta(errorsCommon.CodeInvalidUUID, "the provided list id is not a valid UUID", map[string]any{"field": "listId"}))
		return
	}

	items, err := a.listService.ListItems(r.Context(), listID)
	if err != nil {
		switch {
		case errors.Is(err, middleware.ErrUnauthorized):
			errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
		case errors.Is(err, listService.ErrListNotFound):
			errorsCommon.NotFound(w, errorsCommon.NewError(errorsCommon.CodeListNotFound, "list not found"))
		default:
			errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeDBAccessFailure, "could not list the items"))
		}
		return
	}

	if err := json.NewEncoder(w).Encode(items.ToDto()); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}
