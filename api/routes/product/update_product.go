package productHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"sl-api/api/model/product"
	"sl-api/api/routes/common/errors"
	validatorUtil "sl-api/util/validator"
)

// Update godoc
//
//	@summary        Update product
//	@description    Update product
//	@tags           products
//	@accept         json
//	@produce        json
//	@param          id      path    string  true    "Product ID"
//	@param          body    body    productModel.Form    true    "Product form"
//	@success        200
//	@failure        400 {object}    errorsCommon.ErrorResponse
//	@failure        404 {object}    errorsCommon.ErrorResponse
//	@failure        422 {object}    errorsCommon.ErrorResponse
//	@failure        500 {object}    errorsCommon.ErrorResponse
//	@router         /products/{id} [put]
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewError(errorsCommon.CodeInvalidUUID, "the provided id is not a valid UUID"))
		return
	}

	form := &productModel.Form{}
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

	product := form.ToModel()
	product.ID = id

	rows, err := a.repository.Update(product)
	if err != nil {
		errorsCommon.WriteDBError(w, err, errorsCommon.NewError(errorsCommon.CodeDBUpdateFailure, "could not update the product"))
		return
	}

	if rows == 0 {
		errorsCommon.NotFound(w, errorsCommon.NewError(errorsCommon.CodeProductNotFound, "product not found"))
		return
	}
}
