package user

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	e "sl-api/api/resource/common/err"
	validatorUtil "sl-api/util/validator"
)

type API struct {
	validator  *validator.Validate
	repository *Repository
}

func New(db *gorm.DB, v *validator.Validate) *API {
	return &API{
		repository: NewRepository(db),
		validator:  v,
	}
}

// Register godoc
//
//	@summary		Register user
//	@description	Register user
//	@tags			users
//	@accept			json
//	@product		json
//	@param			body	body	Form	true	"User form"
//	@success		201
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/users [post]
func (a *API) Register(w http.ResponseWriter, r *http.Request) {
	form := &Form{}
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

	newUser := form.ToModel()
	newUser.ID = uuid.New()

	_, err := a.repository.Create(newUser)
	if err != nil {
		e.ServerError(w, e.RespDBDataInsertFailure)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	user, err := a.repository.Read(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	dto := user.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

func (a *API) ReadByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	user, err := a.repository.ReadByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	dto := user.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}
