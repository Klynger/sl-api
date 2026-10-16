package listHandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"sl-api/api/middleware"
	listItemModel "sl-api/api/model/list_item"
	errorsCommon "sl-api/api/routes/common/errors"
	listService "sl-api/api/services/list"
	validatorUtil "sl-api/util/validator"
)

func (a *API) UpdateItemStatus(w http.ResponseWriter, r *http.Request) {
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

	form := &listItemModel.StatusForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
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

	output, err := a.listService.UpdateItemStatus(r.Context(), listService.UpdateItemStatusInput{
		ListID: listID,
		ItemID: itemID,
		Status: listItemModel.Status(form.Status),
	})
	if err != nil {
		switch {
		case errors.Is(err, middleware.ErrUnauthorized):
			errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
		case errors.Is(err, listService.ErrItemNotFound):
			errorsCommon.NotFound(w, errorsCommon.NewError(errorsCommon.CodeItemNotFound, "item not found"))
		default:
			errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeDBUpdateFailure, "could not update the item status"))
		}
		return
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}
