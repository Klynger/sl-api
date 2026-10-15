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

func (a *API) UpdateItem(w http.ResponseWriter, r *http.Request) {
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

	form := &listItemModel.UpdateForm{}
	// An empty body decodes to all-nil fields, which the service rejects as
	// ErrNoFieldsToUpdate, so an EOF is not a decode error; anything else is.
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

	output, err := a.listService.UpdateItem(r.Context(), listService.UpdateItemInput{
		ListID:   listID,
		ItemID:   itemID,
		Quantity: form.Quantity,
		Unit:     form.Unit,
		Note:     form.Note,
	})
	if err != nil {
		switch {
		case errors.Is(err, middleware.ErrUnauthorized):
			errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
		case errors.Is(err, listService.ErrNoFieldsToUpdate):
			errorsCommon.BadRequest(w, errorsCommon.NewError(errorsCommon.CodeNoFieldsToUpdate, "provide at least one field to update"))
		case errors.Is(err, listService.ErrItemNotFound):
			errorsCommon.NotFound(w, errorsCommon.NewError(errorsCommon.CodeItemNotFound, "item not found"))
		default:
			errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeDBUpdateFailure, "could not update the item"))
		}
		return
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}
