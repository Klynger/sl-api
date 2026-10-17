package groupService

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sl-api/api/middleware"
	"sl-api/api/model/invite"
	"sl-api/api/repositories/invite"
	"sl-api/api/repositories/user"
	groupTypes "sl-api/api/types/group"
)

type CreateInviteInput struct {
	GroupID         uuid.UUID
	InvitedUsername string
}

type CreateInviteOutput struct {
	InviteID uuid.UUID `json:"inviteId"`
}

// CreateInvite invites a user (by username) to a group. The caller must be an
// owner of the group.
func (s *GroupService) CreateInvite(ctx context.Context, input CreateInviteInput) (*CreateInviteOutput, error) {
	authedUserID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	userRepo := userRepository.New(s.db)

	invitedUser, err := userRepo.ReadByUsername(input.InvitedUsername)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to look up invited user: %w", err)
	}

	invitedUserID := invitedUser.ID

	if authedUserID == invitedUserID {
		return nil, ErrInviteSelf
	}

	inviteRepo := inviteRepository.New(s.db)

	inviteData, err := inviteRepo.GetCreateInviteData(&groupTypes.GetCreateInviteDataInput{
		Ctx:           ctx,
		GroupID:       input.GroupID,
		InvitedUserID: invitedUserID,
		AuthedUserID:  authedUserID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get create invite data: %w", err)
	}

	if !inviteData.GroupExists {
		return nil, ErrGroupNotFound
	}

	authedUserMember := inviteData.AuthedUserMember

	// Gate on membership before revealing any group state. A non-member is told
	// the group does not exist, so they cannot use this endpoint to probe whether
	// a group exists or whether someone is already invited to / a member of it.
	if authedUserMember == nil {
		return nil, ErrGroupNotFound
	}

	// Only owners can invite.
	if !authedUserMember.IsOwner() {
		return nil, ErrInsufficientPermissions
	}

	// Idempotent: a repeat invite for someone already invited returns the existing
	// pending invite rather than a conflict, so a retrying client converges on the
	// same invite instead of erroring.
	if inviteData.IsAlreadyInvited {
		existing, err := inviteRepo.FindByGroupAndInvitedUser(input.GroupID, invitedUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to load the existing invite: %w", err)
		}

		return &CreateInviteOutput{
			InviteID: existing.ID,
		}, nil
	}

	if inviteData.IsAlreadyMember {
		return nil, ErrAlreadyMember
	}

	invite, err := inviteRepo.Create(&inviteModel.Invite{
		ID:            uuid.New(),
		GroupID:       input.GroupID,
		SenderID:      authedUserID,
		InvitedUserID: invitedUserID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create invite: %w", err)
	}

	return &CreateInviteOutput{
		InviteID: invite.ID,
	}, nil
}
