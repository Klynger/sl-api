package listHandlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"sl-api/api/middleware"
	errorsCommon "sl-api/api/routes/common/errors"
	listService "sl-api/api/services/list"
)

func (a *API) DeleteItem(w http.ResponseWriter, r *http.Request) {
	listID, err := uuid.Parse(chi.URLParam(r, "listId"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewErrorWithMeta(errorsCommon.CodeInvalidUUID, "the provided list id is not a valid UUID", map[string]any{"field": "listId"}))
		return
	}

	itemID, err := uuid.Parse(chi.URLParam(r, "itemId"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewErrorWithMeta(errorsCommon.CodeInvalidUUID, "the provided item id is not a valid UUID", map[string]any{"field": "itemId"}))
		return
	}

	if err := a.listService.DeleteItem(r.Context(), listID, itemID); err != nil {
		switch {
		case errors.Is(err, middleware.ErrUnauthorized):
			errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
		case errors.Is(err, listService.ErrItemNotFound):
			errorsCommon.NotFound(w, errorsCommon.NewError(errorsCommon.CodeItemNotFound, "item not found"))
		default:
			errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeDBAccessFailure, "could not delete the item"))
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
