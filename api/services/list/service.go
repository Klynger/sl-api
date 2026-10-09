package listService

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	groupMemberRepository "sl-api/api/repositories/group_member"
)

// ListService owns the business logic for shopping lists and their items:
// authorization (group membership) and multi-step operations that must stay
// correct across several repositories. Operations are split one-per-file; they
// are all methods on this one struct.
type ListService struct {
	db *gorm.DB
}

func New(db *gorm.DB) *ListService {
	return &ListService{db: db}
}

// ensureGroupMember returns ErrNotAMember unless the user belongs to the group.
// It is the single home for the membership authorization check shared by every
// list operation.
func (s *ListService) ensureGroupMember(userID, groupID uuid.UUID) error {
	memberRepo := groupMemberRepository.New(s.db)

	isMember, err := memberRepo.ExistsByUserAndGroup(userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to check group membership: %w", err)
	}

	if !isMember {
		return ErrNotAMember
	}

	return nil
}
