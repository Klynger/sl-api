package productHandlers

import (
	"encoding/json"
	"net/http"

	"sl-api/api/model/product"
	"sl-api/api/routes/common/errors"
	validatorUtil "sl-api/util/validator"

	uuid "github.com/google/uuid"
)

// Create godoc
//
//	@summary		Create product
//	@description	Create product
//	@tags			products
//	@accept			json
//	@product		json
//	@param			body	body	productModel.Form	true	"Product form"
//	@success		201
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/products [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
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

	newProduct := form.ToModel()
	newProduct.ID = uuid.New()

	_, err := a.repository.Create(newProduct)
	if err != nil {
		errorsCommon.ServerError(w, errorsCommon.RespDBDataInsertFailure)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
