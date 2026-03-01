package group_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sl-api/api/middleware"
	groupModel "sl-api/api/model/group"
	groupMemberModel "sl-api/api/model/group_member"
	"sl-api/api/services/util"
)

type CreateService struct {
	db            *gorm.DB
	newGroupRepo  func(db *gorm.DB) GroupRepository
	newMemberRepo func(db *gorm.DB) GroupMemberRepository
}

func NewCreateService(db *gorm.DB, newGroupRepo func(db *gorm.DB) GroupRepository, newMemberRepo func(db *gorm.DB) GroupMemberRepository) *CreateService {
	return &CreateService{
		db:            db,
		newGroupRepo:  newGroupRepo,
		newMemberRepo: newMemberRepo,
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
		groupRepo := s.newGroupRepo(tx)
		memberRepo := s.newMemberRepo(tx)

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
