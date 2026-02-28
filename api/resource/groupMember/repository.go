package groupMember

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	groupMemberModel "sl-api/api/model/group_member"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(member *groupMemberModel.GroupMember) (*groupMemberModel.GroupMember, error) {
	if err := r.db.Create(member).Error; err != nil {
		return nil, err
	}

	return member, nil
}

func (r *Repository) Read(id uuid.UUID) (*groupMemberModel.GroupMemberDetails, error) {
	var details groupMemberModel.GroupMemberDetails

	err := r.db.Table("group_members").
		Select("group_members.id, group_members.user_id, group_members.group_id, group_members.roles, "+
			"users.user_name, users.last_name, "+
			"groups.name as group_name").
		Joins("INNER JOIN users ON users.id = group_members.user_id AND users.deleted_at IS NULL").
		Joins("INNER JOIN groups ON groups.id = group_members.group_id AND groups.deleted_at IS NULL").
		Where("group_members.id = ? AND group_members.deleted_at IS NULL", id).
		Scan(&details).Error

	if err != nil {
		return nil, err
	}

	return &details, nil
}

func (r *Repository) List() ([]*groupMemberModel.GroupMemberDetails, error) {
	var details []*groupMemberModel.GroupMemberDetails

	err := r.db.Table("group_members").
		Select("group_members.id, group_members.user_id, group_members.group_id, group_members.roles, " +
			"users.user_name, users.last_name, " +
			"groups.name as group_name").
		Joins("INNER JOIN users ON users.id = group_members.user_id AND users.deleted_at IS NULL").
		Joins("INNER JOIN groups ON groups.id = group_members.group_id AND groups.deleted_at IS NULL").
		Where("group_members.deleted_at IS NULL").
		Scan(&details).Error

	if err != nil {
		return nil, err
	}

	return details, nil
}

func (r *Repository) ListByUser(userID uuid.UUID) ([]*groupMemberModel.GroupMemberDetails, error) {
	var details []*groupMemberModel.GroupMemberDetails

	err := r.db.Table("group_members").
		Select("group_members.id, group_members.user_id, group_members.group_id, group_members.roles, "+
			"users.user_name, users.last_name, "+
			"groups.name as group_name").
		Joins("INNER JOIN users ON users.id = group_members.user_id AND users.deleted_at IS NULL").
		Joins("INNER JOIN groups ON groups.id = group_members.group_id AND groups.deleted_at IS NULL").
		Where("group_members.user_id = ? AND group_members.deleted_at IS NULL", userID).
		Scan(&details).Error

	if err != nil {
		return nil, err
	}

	return details, nil
}

func (r *Repository) ListByGroup(groupID uuid.UUID) ([]*groupMemberModel.GroupMemberDetails, error) {
	var details []*groupMemberModel.GroupMemberDetails

	err := r.db.Table("group_members").
		Select("group_members.id, group_members.user_id, group_members.group_id, group_members.roles, "+
			"users.user_name, users.last_name, "+
			"groups.name as group_name").
		Joins("INNER JOIN users ON users.id = group_members.user_id AND users.deleted_at IS NULL").
		Joins("INNER JOIN groups ON groups.id = group_members.group_id AND groups.deleted_at IS NULL").
		Where("group_members.group_id = ? AND group_members.deleted_at IS NULL", groupID).
		Scan(&details).Error

	if err != nil {
		return nil, err
	}

	return details, nil
}

func (r *Repository) FindByUserAndGroup(userID, groupID uuid.UUID) (*groupMemberModel.GroupMemberDetails, error) {
	var details groupMemberModel.GroupMemberDetails

	err := r.db.Table("group_members").
		Select("group_members.id, group_members.user_id, group_members.group_id, group_members.roles, "+
			"users.user_name, users.last_name, "+
			"groups.name as group_name").
		Joins("INNER JOIN users ON users.id = group_members.user_id AND users.deleted_at IS NULL").
		Joins("INNER JOIN groups ON groups.id = group_members.group_id AND groups.deleted_at IS NULL").
		Where("group_members.user_id = ? AND group_members.group_id = ? AND group_members.deleted_at IS NULL", userID, groupID).
		Scan(&details).Error

	if err != nil {
		return nil, err
	}

	return &details, nil
}

func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&groupMemberModel.GroupMember{})

	return result.RowsAffected, result.Error
}
