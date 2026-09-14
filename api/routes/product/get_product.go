package productHandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"sl-api/api/routes/common/errors"
)

// GetProduct godoc
//
//	@summary		Read product
//	@description	Read product
//	@tags			products
//	@accept			json
//	@product		json
//	@param			id	path		string	true	"Product ID"
//	@success		200	{object}	productModel.DTO
//	@failure		400	{object}	errorsCommon.ErrorResponse
//	@failure		404	{object}	errorsCommon.ErrorResponse
//	@failure		500	{object}	errorsCommon.ErrorResponse
//	@router			/products/{id} [get]
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.NewError(errorsCommon.CodeInvalidUUID, "the provided id is not a valid UUID"))
		return
	}

	product, err := a.repository.Read(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorsCommon.NotFound(w, errorsCommon.NewError(errorsCommon.CodeProductNotFound, "product not found"))
			return
		}

		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeDBAccessFailure, "could not read the product"))
		return
	}

	dto := product.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		errorsCommon.ServerError(w, errorsCommon.NewError(errorsCommon.CodeJSONEncodeFailure, "could not encode the response"))
		return
	}
}
