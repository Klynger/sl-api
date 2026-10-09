package listService

import "errors"

// Sentinel errors the handler maps to HTTP status codes. Provisional — extend
// as the operations take shape.
var (
	ErrNotAMember        = errors.New("user is not a member of this group")
	ErrListNotFound      = errors.New("list not found")
	ErrProductNotFound   = errors.New("product not found")
	ErrItemAlreadyOnList = errors.New("product is already on the list")
)
