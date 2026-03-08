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
//	@failure        400 {object}    err.Error
//	@failure        404
//	@failure        500 {object}    err.Error
//	@router         /products/{id} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.RespInvalidURLParamID)
		return
	}

	rows, err := a.repository.Delete(id)
	if err != nil {
		errorsCommon.ServerError(w, errorsCommon.RespDBDdataRemoveFailure)
		return
	}

	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
