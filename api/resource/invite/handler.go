package invite

import (
	"encoding/json"
	"fmt"
	"net/http"

	inviteModel "sl-api/api/model/invite"
	e "sl-api/api/resource/common/err"
	"sl-api/api/resource/group"
	"sl-api/api/resource/groupMember"
	"sl-api/api/resource/user"
	group_service "sl-api/api/services/group"
	validatorUtil "sl-api/util/validator"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type API struct {
	validator           *validator.Validate
	inviteService       *group_service.CreateInviteService
	acceptInviteService *group_service.AcceptInviteService
}

func newGroupRepo(db *gorm.DB) group_service.GroupRepository {
	return group.NewRepository(db)
}

func newGroupMemberRepo(db *gorm.DB) group_service.GroupMemberRepository {
	return groupMember.NewRepository(db)
}

func newInviteRepo(db *gorm.DB) group_service.InviteRepository {
	return NewRepository(db)
}

func newUserRepo(db *gorm.DB) group_service.UserRepository {
	return user.NewRepository(db)
}

func New(db *gorm.DB, validator *validator.Validate) *API {
	return &API{
		validator: validator,
		inviteService: group_service.NewCreateInviteService(
			db,
			newGroupMemberRepo,
			newGroupRepo,
			newInviteRepo,
			newUserRepo,
		),
		acceptInviteService: group_service.NewAcceptInviteService(
			db,
			newInviteRepo,
			newGroupMemberRepo,
		),
	}
}

// Create godoc
//
//	@summary        Create invite
//	@description    Create a new invite for a user to join a group
//	@tags           invites
//	@accept         json
//	@produce        json
//	@param          body    body    group_service.CreateInviteInput    true    "Create invite form"
//	@success        201     {object}    group_service.CreateInviteOutput
//	@failure        400     {object}    e.Error
//	@failure        401     {object}    e.Error
//	@failure        422     {object}    e.Errors
//	@failure        500     {object}    e.Error
//	@router         /invites [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	form := &inviteModel.Form{}

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

	groupID, err := uuid.Parse(chi.URLParam(r, "groupId"))
	if err != nil {
		e.BadRequest(w, e.RespInvalidUUID)
		return
	}

	invitedUserID, err := uuid.Parse(form.InvitedUserID)
	if err != nil {
		e.BadRequest(w, e.RespInvalidUUID)
		return
	}

	input := group_service.CreateInviteInput{
		GroupID:       groupID,
		InvitedUserID: invitedUserID,
	}

	output, err := a.inviteService.Execute(r.Context(), input)

	if err != nil {
		fmt.Println("Error executing invite service:", err)

		e.ServerError(w, e.RespGenericFailure)
		return

		// TODO: Create these errors
		// switch err {
		// case group_service.ErrGroupNotFound:
		// 	e.NotFound(w, e.RespGroupNotFound)
		// case group_service.ErrUserNotFound:
		// 	e.NotFound(w, e.RespUserNotFound)
		// case group_service.ErrAlreadyInvited:
		// 	e.BadRequest(w, e.RespAlreadyInvited)
		// case group_service.ErrAlreadyMember:
		// 	e.BadRequest(w, e.RespAlreadyMember)
		// case group_service.ErrNotGroupOwner:
		// 	e.Unauthorized(w, e.RespNotGroupOwner)
		// default:
		// 	e.ServerError(w, e.RespDBDataInsertFailure)
		// }
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *API) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	groupID, err := uuid.Parse(chi.URLParam(r, "groupId"))
	if err != nil {
		e.BadRequest(w, e.RespInvalidUUID)
		return
	}

	input := group_service.AcceptInviteInput{
		GroupID: groupID,
	}

	output, err := a.acceptInviteService.Execute(r.Context(), input)

	if err != nil {
		fmt.Println("Error executing accept invite service:", err)

		e.ServerError(w, e.RespGenericFailure)
		return

		// TODO: Create the correct errors
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		fmt.Println("Error encoding accept invite output:", err)
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
