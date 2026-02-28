package product

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	uuid "github.com/google/uuid"
	"gorm.io/gorm"

	productModel "sl-api/api/model/product"
	e "sl-api/api/resource/common/err"
	validatorUtil "sl-api/util/validator"
)

type API struct {
	repository *Repository
	validator  *validator.Validate
}

func New(db *gorm.DB, v *validator.Validate) *API {
	return &API{
		repository: NewRepository(db),
		validator:  v,
	}
}

// List godoc
//
//	@summary		List products
//	@description	List products
//	@tags			products
//	@accept			json
//	@product		json
//	@success		200	{object}	[]productModel.DTO
//	@failure		500	{object}	err.Error
//	@router			/products [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	products, err := a.repository.List()
	if err != nil {
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	if len(products) == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(products.ToDto()); err != nil {
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

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
		e.ServerError(w, e.RespJSONDecodeFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			e.ServerError(w, e.RespJSONEncodeFailure)
			return
		}

		e.ValidationErrors(w, respBody)
		return
	}

	newProduct := form.ToModel()
	newProduct.ID = uuid.New()

	_, err := a.repository.Create(newProduct)
	if err != nil {
		e.ServerError(w, e.RespDBDataInsertFailure)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Read godoc
//
//	@summary		Read product
//	@description	Read product
//	@tags			products
//	@accept			json
//	@product		json
//	@param			id	path		string	true	"Product ID"
//	@success		200	{object}	productModel.DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/products/{id} [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	product, err := a.repository.Read(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	dto := product.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

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
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	form := &productModel.Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		e.ServerError(w, e.RespJSONDecodeFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			e.ServerError(w, e.RespJSONEncodeFailure)
			return
		}

		e.ValidationErrors(w, respBody)
		return
	}

	product := form.ToModel()
	product.ID = id

	rows, err := a.repository.Update(product)
	if err != nil {
		e.ServerError(w, e.RespDBDataUpdateFailure)
		return
	}

	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

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
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	rows, err := a.repository.Delete(id)
	if err != nil {
		e.ServerError(w, e.RespDBDdataRemoveFailure)
		return
	}

	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
