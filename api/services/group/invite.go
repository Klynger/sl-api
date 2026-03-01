package group_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sl-api/api/middleware"
	groupModel "sl-api/api/model/group"
	inviteModel "sl-api/api/model/invite"
	userModel "sl-api/api/model/user"
)

type InviteService struct {
	db                 *gorm.DB
	newGroupMemberRepo func(db *gorm.DB) GroupMemberRepository
	newGroupRepo       func(db *gorm.DB) GroupRepository
	newInviteRepo      func(db *gorm.DB) InviteRepository
	newUserRepo        func(db *gorm.DB) UserRepository
}

type InviteInput struct {
	GroupID       uuid.UUID
	InvitedUserID uuid.UUID
}

type InviteOutput struct {
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

func (s *InviteService) Execute(ctx context.Context, input InviteInput) (*InviteOutput, error) {
	authedUserID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if authedUserID == input.InvitedUserID {
		return nil, fmt.Errorf("INVITE_SELF_ERROR")
	}

	groupCh := make(chan groupResult)
	userCh := make(chan userResult)

	go func() {
		groupRepo := s.newGroupRepo(s.db)

		group, err := groupRepo.Read(input.GroupID)
		groupCh <- groupResult{group: group, err: err}
	}()

	go func() {
		userRepo := s.newUserRepo(s.db)

		user, err := userRepo.Read(input.InvitedUserID)
		userCh <- userResult{user: user, err: err}
	}()

	groupRes := <-groupCh
	userRes := <-userCh

	if groupRes.err != nil {
		return nil, fmt.Errorf("GROUP_NOT_FOUND_ERROR")
	}

	if userRes.err != nil {
		return nil, fmt.Errorf("USER_NOT_FOUND_ERROR")
	}

	inviteRepo := s.newInviteRepo(s.db)

	invite, err := inviteRepo.Create(&inviteModel.Invite{
		ID:            uuid.New(),
		GroupID:       input.GroupID,
		SenderID:      authedUserID,
		InvitedUserID: input.InvitedUserID,
	})

	if err != nil {
		return nil, fmt.Errorf("FAILED_TO_CREATE_INVITE_ERROR: %w", err)
	}

	return &InviteOutput{
		InviteID: invite.ID,
	}, nil
}
