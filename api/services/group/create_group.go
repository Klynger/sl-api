package groupService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sl-api/api/middleware"
	"sl-api/api/model/group"
	"sl-api/api/model/group_member"
	"sl-api/api/repositories/group"
	"sl-api/api/repositories/group_member"
	"sl-api/api/services/util"
)

type CreateGroupInput struct {
	Name string
}

type CreateGroupOutput struct {
	GroupID uuid.UUID `json:"groupId"`
	UserID  uuid.UUID `json:"userId"`
}

// CreateGroup creates a group and makes the authenticated user its owner, in one
// transaction.
func (s *GroupService) CreateGroup(ctx context.Context, input CreateGroupInput) (*CreateGroupOutput, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var output *CreateGroupOutput

	err = util.WithTransaction(s.db, func(tx *gorm.DB) error {
		groupRepo := groupRepository.New(tx)
		memberRepo := groupMemberRepository.New(tx)

		newGroup, err := groupRepo.Create(&groupModel.Group{
			ID:   uuid.New(),
			Name: input.Name,
		})

		if err != nil {
			return fmt.Errorf("failed to create group: %w", err)
		}

		newMember, err := memberRepo.Create(&groupMemberModel.GroupMember{
			ID:      uuid.New(),
			GroupID: newGroup.ID,
			UserID:  userID,
			Roles:   []groupMemberModel.Role{groupMemberModel.RoleOwner},
		})

		if err != nil {
			return fmt.Errorf("failed to create group member: %w", err)
		}

		output = &CreateGroupOutput{
			GroupID: newGroup.ID,
			UserID:  newMember.UserID,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}
