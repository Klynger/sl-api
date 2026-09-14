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
		return nil, ErrInviteSelf
	}

	inviteRepo := inviteRepository.New(s.db)

	inviteData, err := inviteRepo.GetCreateInviteData(&groupTypes.GetCreateInviteDataInput{
		Ctx:           ctx,
		GroupID:       input.GroupID,
		InvitedUserID: input.InvitedUserID,
		AuthedUserID:  authedUserID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get create invite data: %w", err)
	}

	if !inviteData.GroupExists {
		return nil, ErrGroupNotFound
	}

	// TODO: Change this to guarantee idempotency. Return the invite instead of an error if already invited
	if inviteData.IsAlreadyInvited {
		return nil, ErrAlreadyInvited
	}

	if inviteData.IsAlreadyMember {
		return nil, ErrAlreadyMember
	}

	authedUserMember := inviteData.AuthedUserMember

	if authedUserMember == nil {
		return nil, ErrNotAMember
	}

	// Only owners can invite
	if !authedUserMember.IsOwner() {
		return nil, ErrInsufficientPermissions
	}

	invite, err := inviteRepo.Create(&inviteModel.Invite{
		ID:            uuid.New(),
		GroupID:       input.GroupID,
		SenderID:      authedUserID,
		InvitedUserID: input.InvitedUserID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create invite: %w", err)
	}

	return &CreateInviteOutput{
		InviteID: invite.ID,
	}, nil
}
