package group_member

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"database/sql/driver"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string
type Roles []Role

const (
	RoleOwner  Role = "owner"
	RoleMember Role = "member"
)

func (r Roles) Value() (driver.Value, error) {
	if len(r) == 0 {
		return "{}", nil
	}

	strs := make([]string, len(r))
	for i, role := range r {
		strs[i] = string(role)
	}

	return "{" + strings.Join(strs, ",") + "}", nil
}

func (r *Roles) Scan(value any) error {
	if value == nil {
		*r = []Role{}
		return nil
	}

	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("failed to scan Roles: %v", value)
		}
		str = string(bytes)
	}

	str = strings.Trim(str, "{}")
	if str == "" {
		*r = []Role{}
		return nil
	}

	parts := strings.Split(str, ",")
	*r = make([]Role, len(parts))
	for i, part := range parts {
		(*r)[i] = Role(strings.Trim(part, `"`))
	}

	return nil
}

type GroupMember struct {
	ID        uuid.UUID `gorm:"primarykey"`
	UserID    uuid.UUID `gorm:"column:user_id"`
	GroupID   uuid.UUID `gorm:"column:group_id"`
	Roles     Roles     `gorm:"type:text[]"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type PublicGroupMember struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	GroupID uuid.UUID
	Roles   Roles
}

func (d *PublicGroupMember) HasRole(role Role) bool {
	return slices.Contains(d.Roles, role)
}

func (d *PublicGroupMember) IsOwner() bool {
	return d.HasRole(RoleOwner)
}

func (GroupMember) TableName() string {
	return "group_members"
}

type GroupMemberDetails struct {
	ID           string
	UserID       string
	UserName     string
	UserLastName string
	GroupID      string
	GroupName    string
	Roles        []string
}

func (d *GroupMemberDetails) ToDto() *DTO {
	return &DTO{
		ID:           d.ID,
		UserID:       d.UserID,
		UserName:     d.UserName,
		UserLastName: d.UserLastName,
		GroupID:      d.GroupID,
		GroupName:    d.GroupName,
		Roles:        d.Roles,
	}
}

func (d *GroupMemberDetails) HasRole(role Role) bool {
	return slices.Contains(d.Roles, string(role))
}

func (d *GroupMemberDetails) IsOwner() bool {
	return d.HasRole(RoleOwner)
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
