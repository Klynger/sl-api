package userRepository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"sl-api/api/model/group_member"
	"sl-api/api/model/user"
)

type UserRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *userModel.User) (*userModel.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Read(id uuid.UUID) (*userModel.User, error) {
	user := &userModel.User{}
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) ReadByUsername(username string) (*userModel.UserCredentials, error) {
	user := &userModel.UserCredentials{}
	if err := r.db.Select("id", "username", "password").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) DeleteWithCascade(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", id).Delete(&groupMemberModel.GroupMember{}).Error; err != nil {
			return err
		}

		return tx.Delete(&userModel.User{}, id).Error
	})
}
