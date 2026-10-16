package groupService

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"sl-api/api/middleware"
	groupMemberModel "sl-api/api/model/group_member"
	groupMemberRepository "sl-api/api/repositories/group_member"
)

// ListGroupMembers returns the members of a group, each with their roles. The
// caller must be a member; a non-member is reported as a missing group, so the
// endpoint never reveals a group the caller cannot see.
func (s *GroupService) ListGroupMembers(ctx context.Context, groupID uuid.UUID) (groupMemberModel.GroupMembersDetails, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	memberRepo := groupMemberRepository.New(s.db)

	isMember, err := memberRepo.ExistsByUserAndGroup(userID, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to check group membership: %w", err)
	}

	if !isMember {
		return nil, ErrGroupNotFound
	}

	members, err := memberRepo.ListByGroup(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list the group's members: %w", err)
	}

	return groupMemberModel.GroupMembersDetails(members), nil
}
