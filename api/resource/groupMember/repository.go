package groupMember

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(member *GroupMember) (*GroupMember, error) {
	if err := r.db.Create(member).Error; err != nil {
		return nil, err
	}

	return member, nil
}

func (r *Repository) Read(id uuid.UUID) (*GroupMember, error) {
	member := &GroupMember{}

	if err := r.db.Where("id = ?", id).First(&member).Error; err != nil {
		return nil, err
	}

	return member, nil
}

func (r *Repository) List() {

}

func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&GroupMember{})

	return result.RowsAffected, result.Error
}
