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
//	@failure        400 {object}    err.Error
//	@failure        404
//	@failure        422 {object}    err.Errors
//	@failure        500 {object}    err.Error
//	@router         /products/{id} [put]
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.RespInvalidURLParamID)
		return
	}

	form := &productModel.Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		errorsCommon.ServerError(w, errorsCommon.RespJSONDecodeFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			errorsCommon.ServerError(w, errorsCommon.RespJSONEncodeFailure)
			return
		}

		errorsCommon.ValidationErrors(w, respBody)
		return
	}

	product := form.ToModel()
	product.ID = id

	rows, err := a.repository.Update(product)
	if err != nil {
		errorsCommon.ServerError(w, errorsCommon.RespDBDataUpdateFailure)
		return
	}

	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
