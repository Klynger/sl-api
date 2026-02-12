package user

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
	"net/http"
	e "sl-api/api/resource/common/err"
	validatorUtil "sl-api/util/validator"
)

type API struct {
	validator        *validator.Validate
	repository       *Repository
	authSessionStore *sessions.CookieStore
}

func New(db *gorm.DB, authSessionStore *sessions.CookieStore, v *validator.Validate) *API {
	return &API{
		repository:       NewRepository(db),
		validator:        v,
		authSessionStore: authSessionStore,
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

func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	form := &LoginForm{}
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

	userCredentials, err := a.repository.ReadByUsername(form.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// TODO: Do something to avoid timing attacks
			w.WriteHeader(http.StatusNotFound)
			return
		}

		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	validationResult := userCredentials.ValidateCredentials(form)

	if !validationResult {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session, err := a.authSessionStore.Get(r, "auth-session")
	if err != nil {
		e.ServerError(w, e.RespSessionAccessFailure)
		return
	}

	session.Values["user_id"] = userCredentials.ID.String()
	session.Values["username"] = userCredentials.Username

	session.Options = &sessions.Options{
		MaxAge:   3600,  // 1 hour
		HttpOnly: false, // SECURITY FLAW: Allows JavaScript access
		Secure:   false, // SECURITY FLAW: Allows HTTP access (not just HTTPS)
	}

	err = session.Save(r, w)
	if err != nil {
		fmt.Println("Session save error: ", err)
		e.ServerError(w, e.RespGenericFailure)
		return
	}

	resp := &LoginResponse{
		Message:  "Login successful",
		Username: userCredentials.Username,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type LoginResponse struct {
	Message  string    `json:"message"`
	Username string    `json:"username"`
}
