package listHandlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"sl-api/api/middleware"
	listItemModel "sl-api/api/model/list_item"
	errorsCommon "sl-api/api/routes/common/errors"
	listService "sl-api/api/services/list"
	validatorUtil "sl-api/util/validator"
)

func (a *API) IncreaseItemQuantity(w http.ResponseWriter, r *http.Request) {
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

	form := &listItemModel.IncrementForm{}
	// An empty body is allowed and means "increment by the default of 1", so an
	// EOF from the decoder is not an error; anything else malformed is.
	if err := json.NewDecoder(r.Body).Decode(form); err != nil && !errors.Is(err, io.EOF) {
		errorsCommon.BadRequest(w, errorsCommon.NewError(errorsCommon.CodeInvalidJSON, "invalid JSON in request body"))
		return
	}

	if err := a.validator.Struct(form); err != nil {
		items := validatorUtil.ToErrResponse(err)
		if items == nil {
			errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeInternalError, "unexpected validation error"))
			return
		}

		errorsCommon.ValidationErrors(w, items...)
		return
	}

	by := 1
	if form.By != nil {
		by = *form.By
	}

	output, err := a.listService.IncreaseQuantity(r.Context(), listService.IncreaseQuantityInput{
		ListID: listID,
		ItemID: itemID,
		By:     by,
	})
	if err != nil {
		switch {
		case errors.Is(err, middleware.ErrUnauthorized):
			errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
		case errors.Is(err, listService.ErrItemNotFound):
			errorsCommon.NotFound(w, errorsCommon.NewError(errorsCommon.CodeItemNotFound, "item not found"))
		default:
			errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeDBUpdateFailure, "could not increase the item quantity"))
		}
		return
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}
