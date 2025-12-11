package product

import "net/http"

type API struct{}

// Create godoc
//
//	@summary		Create product
//	@description	Create product
//	@tags			products
//	@accept			json
//	@product		json
//	@param			body	body	Form	true	"Product form"
//	@success		201
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/products [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {}

// Read godoc
//
//	@summary		Read product
//	@description	Read product
//	@tags			products
//	@accept			json
//	@product		json
//	@param			id	path		string	true	"Product ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/products/{id} [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {}
