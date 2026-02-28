package create_group

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sl-api/api/middleware"
	groupModel "sl-api/api/model/group"
	groupMemberModel "sl-api/api/model/group_member"
)

type GroupRepository interface {
	Create(group *groupModel.Group) (*groupModel.Group, error)
}

type GroupMemberRepository interface {
	Create(member *groupMemberModel.GroupMember) (*groupMemberModel.GroupMember, error)
}

type Service struct {
	db            *gorm.DB
	newGroupRepo  func(db *gorm.DB) GroupRepository
	newMemberRepo func(db *gorm.DB) GroupMemberRepository
}

func NewService(db *gorm.DB, newGroupRepo func(db *gorm.DB) GroupRepository, newMemberRepo func(db *gorm.DB) GroupMemberRepository) *Service {
	return &Service{
		db:            db,
		newGroupRepo:  newGroupRepo,
		newMemberRepo: newMemberRepo,
	}
}

type ExecuteInput struct {
	Name string
}

type ExecuteOutput struct {
	GroupID uuid.UUID
	UserID  uuid.UUID
}

func (s *Service) Execute(ctx context.Context, input ExecuteInput) (*ExecuteOutput, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var output *ExecuteOutput

	err = s.WithTransaction(ctx, func(tx *gorm.DB) error {
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

		output = &ExecuteOutput{
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

func (s *Service) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := s.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
