package listService

import "errors"

// Sentinel errors for list operations; handlers map these to HTTP responses.
//
// Access control is non-leaking: an action the caller may not see returns the
// same not-found error a missing resource would, so membership is never
// revealed. ensureGroupMember yields ErrNotAMember internally, and each
// operation translates it into the not-found error that fits its scope (group,
// list, or item) before returning.
var (
	ErrNotAMember           = errors.New("user is not a member of this group")
	ErrGroupNotFound        = errors.New("group not found")
	ErrListNotFound         = errors.New("list not found")
	ErrItemNotFound         = errors.New("item not found")
	ErrProductNotFound      = errors.New("product not found")
	ErrItemAlreadyOnList    = errors.New("product is already on the list")
	ErrQuantityBelowMinimum = errors.New("quantity cannot drop below 1")
	ErrNoFieldsToUpdate     = errors.New("no fields to update")
)
