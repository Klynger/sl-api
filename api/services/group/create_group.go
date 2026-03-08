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

type CreateService struct {
	db *gorm.DB
}

func NewCreateService(db *gorm.DB) *CreateService {
	return &CreateService{
		db: db,
	}
}

type CreateInput struct {
	Name string
}

type CreateOutput struct {
	GroupID uuid.UUID
	UserID  uuid.UUID
}

func (s *CreateService) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var output *CreateOutput

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

		output = &CreateOutput{
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
