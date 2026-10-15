package groupService

import (
	"gorm.io/gorm"
)

// GroupService owns the business logic for groups and their invites:
// authorization (group membership) and multi-step operations that must stay
// correct across several repositories. Operations are split one-per-file; they
// are all methods on this one struct.
type GroupService struct {
	db *gorm.DB
}

func New(db *gorm.DB) *GroupService {
	return &GroupService{db: db}
}
