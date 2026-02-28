package group

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sl-api/api/resource/groupMember"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(group *Group) (*Group, error) {
	if err := r.db.Create(group).Error; err != nil {
		return nil, err
	}

	return group, nil
}

func (r *Repository) Read(id uuid.UUID) (*Group, error) {
	group := &Group{}
	if err := r.db.Where("id = ?", id).First(&group).Error; err != nil {
		return nil, err
	}

	return group, nil
}

func (r *Repository) DeleteWithCascade(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ?", id).Delete(&groupMember.GroupMember{}).Error; err != nil {
			return err
		}

		return tx.Delete(&Group{}, id).Error
	})
}
