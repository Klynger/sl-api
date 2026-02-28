package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	groupMemberModel "sl-api/api/model/group_member"
	userModel "sl-api/api/model/user"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(user *userModel.User) (*userModel.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) Read(id uuid.UUID) (*userModel.User, error) {
	user := &userModel.User{}
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) ReadByUsername(username string) (*userModel.UserCredentials, error) {
	user := &userModel.UserCredentials{}
	if err := r.db.Select("id", "username", "password").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) DeleteWithCascade(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", id).Delete(&groupMemberModel.GroupMember{}).Error; err != nil {
			return err
		}

		return tx.Delete(&userModel.User{}, id).Error
	})
}
