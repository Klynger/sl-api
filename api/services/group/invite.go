package groupService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sl-api/api/middleware"
	 "sl-api/api/model/group"
	 "sl-api/api/model/invite"
	 "sl-api/api/model/user"
	"sl-api/api/repositories/invite"
	groupTypes "sl-api/api/types/group"
)

type CreateInviteService struct {
	db *gorm.DB
}

func NewCreateInviteService(
	db *gorm.DB,
) *CreateInviteService {
	return &CreateInviteService{
		db: db,
	}
}

type CreateInviteInput struct {
	GroupID       uuid.UUID
	InvitedUserID uuid.UUID
}

type CreateInviteOutput struct {
	InviteID uuid.UUID
}

type groupResult struct {
	group *groupModel.Group
	err   error
}

type userResult struct {
	user *userModel.User
	err  error
}

func (s *CreateInviteService) Execute(ctx context.Context, input CreateInviteInput) (*CreateInviteOutput, error) {
	authedUserID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if authedUserID == input.InvitedUserID {
		return nil, fmt.Errorf("INVITE_SELF_ERROR")
	}

	inviteRepo := inviteRepository.New(s.db)

	inviteData, err := inviteRepo.GetCreateInviteData(&groupTypes.GetCreateInviteDataInput{
		Ctx:           ctx,
		GroupID:       input.GroupID,
		InvitedUserID: input.InvitedUserID,
		AuthedUserID:  authedUserID,
	})

	if err != nil {
		return nil, fmt.Errorf("FAILED_TO_GET_CREATE_INVITE_DATA_ERROR: %w", err)
	}

	if !inviteData.GroupExists {
		return nil, fmt.Errorf("GROUP_NOT_FOUND_ERROR")
	}

	// TODO: Change this to guarantee idempotency. Return the invite instead of an error if already invited
	if inviteData.IsAlreadyInvited {
		return nil, fmt.Errorf("ALREADY_INVITED_ERROR")
	}

	if inviteData.IsAlreadyMember {
		return nil, fmt.Errorf("ALREADY_MEMBER_ERROR")
	}

	authedUserMember := inviteData.AuthedUserMember

	if authedUserMember == nil {
		return nil, fmt.Errorf("NOT_A_MEMBER_ERROR")
	}

	// Only owners can invite
	if !authedUserMember.IsOwner() {
		return nil, fmt.Errorf("INSUFFICIENT_PERMISSIONS_ERROR")
	}

	invite, err := inviteRepo.Create(&inviteModel.Invite{
		ID:            uuid.New(),
		GroupID:       input.GroupID,
		SenderID:      authedUserID,
		InvitedUserID: input.InvitedUserID,
	})

	if err != nil {
		return nil, fmt.Errorf("FAILED_TO_CREATE_INVITE_ERROR: %w", err)
	}

	return &CreateInviteOutput{
		InviteID: invite.ID,
	}, nil
}
