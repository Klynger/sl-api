package group

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	groupModel "sl-api/api/model/group"
	e "sl-api/api/resource/common/err"
	"sl-api/api/resource/groupMember"
	"sl-api/api/services/group"
	validatorUtil "sl-api/util/validator"
)

type API struct {
	validator          *validator.Validate
	createGroupService *group_service.CreateService
}

func newGroupRepo(db *gorm.DB) group_service.GroupRepository {
	return NewRepository(db)
}

func newMemberRepo(db *gorm.DB) group_service.GroupMemberRepository {
	return groupMember.NewRepository(db)
}

func New(db *gorm.DB, validator *validator.Validate) *API {
	return &API{
		validator:          validator,
		createGroupService: group_service.NewCreateService(db, newGroupRepo, newMemberRepo),
	}
}

// Create godoc
//
//	@summary        Create group
//	@description    Create a new group with the authenticated user as owner
//	@tags           groups
//	@accept         json
//	@produce        json
//	@param          body    body    groupModel.Form    true    "Group form"
//	@success        201     {object}    groupModel.DTO
//	@failure        400     {object}    err.Error
//	@failure        401     {object}    err.Error
//	@failure        422     {object}    err.Errors
//	@failure        500     {object}    err.Error
//	@router         /groups [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	form := &groupModel.Form{}
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

	input := group_service.CreateInput{
		Name: form.Name,
	}

	output, err := a.createGroupService.Execute(r.Context(), input)
	if err != nil {
		e.ServerError(w, e.RespDBDataInsertFailure)
		return
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
