package invite

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	groupMemberModel "sl-api/api/model/group_member"
	inviteModel "sl-api/api/model/invite"
	userModel "sl-api/api/model/user"
	group_service "sl-api/api/services/group"
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

func (r *Repository) GetCreateInviteData(input *group_service.GetCreateInviteDataInput) (*group_service.GetCreateInviteDataOutput, error) {
	type queryResult struct {
		// Flags
		GroupExists      bool
		IsAlreadyInvited bool
		IsAlreadyMember  bool
		// Authed user fields (nullable)
		AuthedUserID   *uuid.UUID
		AuthedName     *string
		AuthedLastName *string
		// Authed member (nullable)
		AuthedMemberID      *uuid.UUID
		AuthedMemberUserID  *uuid.UUID
		AuthedMemberGroupID *uuid.UUID
		AuthedMemberRoles   *groupMemberModel.Roles
		// Invited user fields (nullable)
		InvitedUserID   *uuid.UUID
		InvitedName     *string
		InvitedLastName *string
	}

	var result queryResult

	query := `
		SELECT
			-- Check if group exists
			CASE WHEN g.id IS NOT NULL THEN true ELSE false END as group_exists,
			
			-- Check if already invited
			CASE WHEN i.id IS NOT NULL THEN true ELSE false END as is_already_invited,
			
			-- Check if already a member (invited user)
			CASE WHEN gm.id IS NOT NULL THEN true ELSE false END as is_already_member,
			
			-- Authed user data
			au.id as authed_user_id,
			au.user_name as authed_name,
			au.last_name as authed_last_name,
			
			-- Invited user data
			iu.id as invited_user_id,
			iu.user_name as invited_name,
			iu.last_name as invited_last_name,
			
		-- Authed user's member data
		au_gm.id as authed_member_id,
		au_gm.user_id as authed_member_user_id,
		au_gm.group_id as authed_member_group_id,
		au_gm.roles as authed_member_roles
		
	FROM users AS au
		LEFT JOIN users AS iu 
			ON iu.id = ? 
			AND iu.deleted_at IS NULL
		LEFT JOIN groups AS g 
			ON g.id = ? 
			AND g.deleted_at IS NULL
		LEFT JOIN invites AS i 
			ON i.group_id = ? 
			AND i.invited_user_id = ? 
			AND i.deleted_at IS NULL
		LEFT JOIN group_members AS gm 
			ON gm.group_id = ? 
			AND gm.user_id = ? 
			AND gm.deleted_at IS NULL
		LEFT JOIN group_members AS au_gm
			ON au_gm.group_id = ?
			AND au_gm.user_id = ?
			AND au_gm.deleted_at IS NULL
			
		WHERE au.id = ? 
			AND au.deleted_at IS NULL
	`

	err := r.db.WithContext(input.Ctx).Raw(
		query,
		input.InvitedUserID, // iu.id
		input.GroupID,       // g.id
		input.GroupID,       // i.group_id
		input.InvitedUserID, // i.invited_user_id
		input.GroupID,       // gm.group_id
		input.InvitedUserID, // gm.user_id (invited user's membership)
		input.GroupID,       // au_gm.group_id
		input.AuthedUserID,  // au_gm.user_id (authed user's membership)
		input.AuthedUserID,  // au.id

	).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	inviteData := &group_service.GetCreateInviteDataOutput{
		GroupExists:      result.GroupExists,
		IsAlreadyInvited: result.IsAlreadyInvited,
		IsAlreadyMember:  result.IsAlreadyMember,
	}

	if result.AuthedUserID != nil {
		inviteData.AuthedUser = &userModel.PublicUser{
			ID:       *result.AuthedUserID,
			Name:     *result.AuthedName,
			LastName: *result.AuthedLastName,
		}
	}

	if result.InvitedUserID != nil {
		inviteData.InvitedUser = &userModel.PublicUser{
			ID:       *result.InvitedUserID,
			Name:     *result.InvitedName,
			LastName: *result.InvitedLastName,
		}
	}

	if result.AuthedMemberID != nil {
		inviteData.AuthedUserMember = &groupMemberModel.PublicGroupMember{
			ID:      *result.AuthedMemberID,
			UserID:  *result.AuthedMemberUserID,
			GroupID: *result.AuthedMemberGroupID,
			Roles:   *result.AuthedMemberRoles,
		}
	}

	return inviteData, nil
}

func (r *Repository) DeleteAllInvitesFromUser(invitedUserID uuid.UUID) (int64, error) {
	result := r.db.Where("invited_user_id = ?", invitedUserID).Delete(&inviteModel.Invite{})

	return result.RowsAffected, result.Error
}

func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&inviteModel.Invite{})

	return result.RowsAffected, result.Error
}
