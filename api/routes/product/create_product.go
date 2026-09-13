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
//	@failure		400	{object}	errorsCommon.ErrorResponse
//	@failure		422	{object}	errorsCommon.ErrorResponse
//	@failure		500	{object}	errorsCommon.ErrorResponse
//	@router			/products [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
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

	newProduct := form.ToModel()
	newProduct.ID = uuid.New()

	_, err := a.repository.Create(newProduct)
	if err != nil {
		errorsCommon.WriteDBError(w, err, errorsCommon.NewError(errorsCommon.CodeDBInsertFailure, "could not save the product"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
