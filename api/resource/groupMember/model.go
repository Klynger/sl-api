package groupMember

import (
	"slices"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleOwner  Role = "owner"
	RoleMember Role = "member"
)

type GroupMember struct {
	ID        uuid.UUID `gorm:"primarykey"`
	UserID    uuid.UUID `gorm:"column:user_id"`
	GroupID   uuid.UUID `gorm:"column:group_id"`
	Roles     []Role    `gorm:"type:text[]"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (GroupMember) TableName() string {
	return "group_members"
}

type DTO struct {
	ID           string   `json:"id"`
	UserID       string   `json:"userId"`
	UserName     string   `json:"userName"`
	UserLastName string   `json:"userLastName"`
	GroupID      string   `json:"groupId"`
	GroupName    string   `json:"groupName"`
	Roles        []string `json:"roles"`
}

type Form struct {
	UserID  string   `json:"userId" validate:"required,uuid"`
	GroupID string   `json:"groupId" validate:"required,uuid"`
	Roles   []string `json:"roles" validate:"required,dive,oneof=owner member"`
}

func (m *GroupMember) ToDto() *DTO {
	roles := make([]string, len(m.Roles))
	for i, role := range m.Roles {
		roles[i] = string(role)
	}

	return &DTO{
		ID:      m.ID.String(),
		UserID:  m.UserID.String(),
		GroupID: m.GroupID.String(),
		Roles:   roles,
	}
}

func (m *GroupMember) HasRole(role Role) bool {
	return slices.Contains(m.Roles, role)
}

func (m *GroupMember) IsOwner() bool {
	return m.HasRole(RoleOwner)
}

func (f *Form) ToModel() (*GroupMember, error) {
	userID, err := uuid.Parse(f.UserID)
	if err != nil {
		return nil, err
	}

	rolesStr := f.Roles
	if len(rolesStr) == 0 {
		rolesStr = []string{string(RoleMember)}
	}

	memberRoles := make([]Role, len(rolesStr))
	for i, role := range rolesStr {
		memberRoles[i] = Role(role)
	}

	groupID, err := uuid.Parse(f.GroupID)
	if err != nil {
		return nil, err
	}

	return &GroupMember{
		UserID:  userID,
		GroupID: groupID,
		Roles:   memberRoles,
	}, nil
}

type GroupMembers []*GroupMember

func (ms GroupMembers) ToDto() []*DTO {
	dtos := make([]*DTO, len(ms))
	for i, v := range ms {
		dtos[i] = v.ToDto()
	}

	return dtos
}
