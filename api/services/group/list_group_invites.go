package groupService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sl-api/api/middleware"
	inviteModel "sl-api/api/model/invite"
	groupMemberRepository "sl-api/api/repositories/group_member"
	inviteRepository "sl-api/api/repositories/invite"
)

type ListGroupInvitesService struct {
	db *gorm.DB
}

func NewListGroupInvitesService(
	db *gorm.DB,
) *ListGroupInvitesService {
	return &ListGroupInvitesService{
		db: db,
	}
}

// Execute returns the pending invites for a group. The caller must be a member;
// a non-member is reported as a missing group, so the endpoint never reveals a
// group the caller cannot see.
func (s *ListGroupInvitesService) Execute(ctx context.Context, groupID uuid.UUID) (inviteModel.PendingInvites, error) {
	authedUserID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	memberRepo := groupMemberRepository.New(s.db)

	isMember, err := memberRepo.ExistsByUserAndGroup(authedUserID, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to check group membership: %w", err)
	}

	if !isMember {
		return nil, ErrGroupNotFound
	}

	inviteRepo := inviteRepository.New(s.db)

	invites, err := inviteRepo.ListPendingForGroup(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list the group's pending invites: %w", err)
	}

	return invites, nil
}
