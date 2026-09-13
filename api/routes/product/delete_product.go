package productHandlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"sl-api/api/routes/common/errors"
)

// Delete godoc
//
//	@summary        Delete product
//	@description    Delete product
//	@tags           products
//	@accept         json
//	@produce        json
//	@param          id  path    string  true    "Product ID"
//	@success        200
//	@failure        400 {object}    errorsCommon.ErrorResponse
//	@failure        404 {object}    errorsCommon.ErrorResponse
//	@failure        500 {object}    errorsCommon.ErrorResponse
//	@router         /products/{id} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewError(errorsCommon.CodeInvalidUUID, "the provided id is not a valid UUID"))
		return
	}

	rows, err := a.repository.Delete(id)
	if err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeDBRemoveFailure, "could not delete the product"))
		return
	}

	if rows == 0 {
		errorsCommon.NotFound(w, errorsCommon.NewError(errorsCommon.CodeProductNotFound, "product not found"))
		return
	}
}
