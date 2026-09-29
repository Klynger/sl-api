package listRepository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	listModel "sl-api/api/model/list"
)

type ListRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *ListRepository {
	return &ListRepository{
		db: db,
	}
}

func (r *ListRepository) Create(list *listModel.List) (*listModel.List, error) {
	if err := r.db.Create(list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r *ListRepository) Read(id uuid.UUID) (*listModel.List, error) {
	list := &listModel.List{}
	if err := r.db.Where("id = ?", id).First(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r *ListRepository) ListByGroup(groupID uuid.UUID) (listModel.Lists, error) {
	lists := make([]*listModel.List, 0)

	if err := r.db.Where("group_id = ?", groupID).Find(&lists).Error; err != nil {
		return nil, err
	}

	return lists, nil
}

func (r *ListRepository) Update(list *listModel.List) (int64, error) {
	result := r.db.Model(&listModel.List{}).
		Select("Name", "Status", "UpdatedAt").
		Where("id = ?", list.ID).
		Updates(list)

	return result.RowsAffected, result.Error
}

func (r *ListRepository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&listModel.List{})

	return result.RowsAffected, result.Error
}
