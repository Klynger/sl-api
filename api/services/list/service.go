package listService

import "gorm.io/gorm"

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

// TODO: shared group-membership authorization helper.
// Where this lives (a method here, the handler, or middleware) is the decision
// we're about to make. Every operation below will gate on it before touching
// list or item data.
