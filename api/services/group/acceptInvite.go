package group_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sl-api/api/middleware"
	groupMemberModel "sl-api/api/model/group_member"
	"sl-api/api/services/util"
)

type AcceptInviteService struct {
	db                 *gorm.DB
	newInviteRepo      newInviteRepoType
	newGroupMemberRepo newGroupMemberRepoType
}

func NewAcceptInviteService(
	db *gorm.DB,
	newInviteRepo newInviteRepoType,
	newGroupMemberRepo newGroupMemberRepoType,
) *AcceptInviteService {
	return &AcceptInviteService{
		db:                 db,
		newInviteRepo:      newInviteRepo,
		newGroupMemberRepo: newGroupMemberRepo,
	}
}

type AcceptInviteInput struct {
	GroupID uuid.UUID
}

type AcceptInviteOutput struct {
	GroupMember groupMemberModel.GroupMember
}

func (s *AcceptInviteService) Execute(ctx context.Context, input AcceptInviteInput) (*AcceptInviteOutput, error) {
	authedUserID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var output *AcceptInviteOutput

	err = util.WithTransaction(s.db, func(tx *gorm.DB) error {
		inviteRepo := s.newInviteRepo(tx)
		groupMemberRepo := s.newGroupMemberRepo(tx)

		inviteData, err := inviteRepo.GetAcceptInviteData(input.GroupID, authedUserID)

		if err != nil {
			return fmt.Errorf("FAILED_TO_GET_INVITE")
		}

		memberRoles := []groupMemberModel.Role{groupMemberModel.RoleMember}

		groupMember := groupMemberModel.GroupMember{
			ID:      uuid.New(),
			UserID:  authedUserID,
			GroupID: input.GroupID,
			Roles:   memberRoles,
		}

		groupMemberRepo.Create(&groupMember)

		if err := inviteRepo.Delete(inviteData.ID); err != nil {
			return fmt.Errorf("FAILED_TO_DELETE_INVITE")
		}

		output = &AcceptInviteOutput{
			GroupMember: groupMember,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}
