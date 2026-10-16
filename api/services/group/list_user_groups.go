package groupService

import (
	"context"
	"fmt"

	"sl-api/api/middleware"
	groupMemberModel "sl-api/api/model/group_member"
	groupMemberRepository "sl-api/api/repositories/group_member"
)

// ListUserGroups returns the groups the authenticated caller belongs to, with
// the caller's role in each. The result is naturally scoped to that user via
// their memberships, so there is nothing to leak.
func (s *GroupService) ListUserGroups(ctx context.Context) (groupMemberModel.GroupMembersDetails, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	memberRepo := groupMemberRepository.New(s.db)
	memberships, err := memberRepo.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list the user's groups: %w", err)
	}

	return groupMemberModel.GroupMembersDetails(memberships), nil
}
