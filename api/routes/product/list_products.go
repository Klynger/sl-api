package productHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	errorsCommon "sl-api/api/routes/common/errors"
)

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
		errorsCommon.ServerError(w, errorsCommon.RespDBDataAccessFailure)
		return
	}

	if len(products) == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(products.ToDto()); err != nil {
		errorsCommon.ServerError(w, errorsCommon.RespJSONEncodeFailure)
		return
	}
}
