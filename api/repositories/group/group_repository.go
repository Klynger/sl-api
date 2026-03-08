package groupRepository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	groupModel "sl-api/api/model/group"
	groupMemberModel "sl-api/api/model/group_member"
)

type GroupRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *GroupRepository {
	return &GroupRepository{
		db: db,
	}
}

func (r *GroupRepository) Create(group *groupModel.Group) (*groupModel.Group, error) {
	if err := r.db.Create(group).Error; err != nil {
		return nil, err
	}

	return group, nil
}

func (r *GroupRepository) Read(id uuid.UUID) (*groupModel.Group, error) {
	group := &groupModel.Group{}
	if err := r.db.Where("id = ?", id).First(&group).Error; err != nil {
		return nil, err
	}

	return group, nil
}

func (r *GroupRepository) DeleteWithCascade(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ?", id).Delete(&groupMemberModel.GroupMember{}).Error; err != nil {
			return err
		}

		return tx.Delete(&groupModel.Group{}, id).Error
	})
}
