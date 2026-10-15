package groupService

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"sl-api/api/middleware"
	inviteModel "sl-api/api/model/invite"
	inviteRepository "sl-api/api/repositories/invite"
)

type ListUserInvitesService struct {
	db *gorm.DB
}

func NewListUserInvitesService(
	db *gorm.DB,
) *ListUserInvitesService {
	return &ListUserInvitesService{
		db: db,
	}
}

// Execute returns the pending invites addressed to the authenticated caller.
// The result is naturally scoped to that user, so there is nothing to leak.
func (s *ListUserInvitesService) Execute(ctx context.Context) (inviteModel.PendingInvites, error) {
	authedUserID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	inviteRepo := inviteRepository.New(s.db)

	invites, err := inviteRepo.ListPendingForUser(authedUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to list the user's pending invites: %w", err)
	}

	return invites, nil
}
