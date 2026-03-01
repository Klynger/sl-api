package invite

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	inviteModel "sl-api/api/model/invite"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(invite *inviteModel.Invite) (*inviteModel.Invite, error) {
	if err := r.db.Create(invite).Error; err != nil {
		return nil, err
	}

	return invite, nil
}

func (r *Repository) DeleteAllInvitesFromUser(invitedUserID uuid.UUID) (int64, error) {
	result := r.db.Where("invited_user_id = ?", invitedUserID).Delete(&inviteModel.Invite{})

	return result.RowsAffected, result.Error
}

func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&inviteModel.Invite{})

	return result.RowsAffected, result.Error
}
