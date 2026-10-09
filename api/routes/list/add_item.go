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

func (a *API) AddItem(w http.ResponseWriter, r *http.Request) {
	listID, err := uuid.Parse(chi.URLParam(r, "listId"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewErrorWithMeta(errorsCommon.CodeInvalidUUID, "the provided list id is not a valid UUID", map[string]any{"field": "listId"}))
		return
	}

	form := &listItemModel.Form{}
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

	productID, err := uuid.Parse(form.ProductID)
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewErrorWithMeta(errorsCommon.CodeInvalidUUID, "the provided product id is not a valid UUID", map[string]any{"field": "productId"}))
		return
	}

	output, err := a.listService.AddItem(r.Context(), listService.AddItemInput{
		ListID:    listID,
		ProductID: productID,
		Quantity:  form.Quantity,
		Unit:      form.Unit,
		Note:      form.Note,
	})
	if err != nil {
		switch {
		case errors.Is(err, middleware.ErrUnauthorized):
			errorsCommon.Unauthorized(w, errorsCommon.NewError(errorsCommon.CodeUnauthorized, "authentication required"))
		case errors.Is(err, listService.ErrNotAMember):
			errorsCommon.Forbidden(w, errorsCommon.NewError(errorsCommon.CodeNotAMember, "you are not a member of this group"))
		case errors.Is(err, listService.ErrListNotFound):
			errorsCommon.NotFound(w, errorsCommon.NewError(errorsCommon.CodeListNotFound, "list not found"))
		case errors.Is(err, listService.ErrProductNotFound):
			errorsCommon.NotFound(w, errorsCommon.NewError(errorsCommon.CodeProductNotFound, "product not found"))
		case errors.Is(err, listService.ErrItemAlreadyOnList):
			errorsCommon.Conflict(w, errorsCommon.NewError(errorsCommon.CodeItemAlreadyOnList, "product is already on the list"))
		default:
			errorsCommon.WriteDBError(w, err, errorsCommon.NewError(errorsCommon.CodeDBInsertFailure, "could not add the item"))
		}
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(output); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}
