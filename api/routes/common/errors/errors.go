package errorsCommon

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	CodeInvalidJSON       = "INVALID_JSON"
	CodeInvalidUUID       = "INVALID_UUID"
	CodeJSONEncodeFailure = "JSON_ENCODE_FAILURE"
	CodeInternalError     = "INTERNAL_ERROR"

	CodeDBAccessFailure = "DB_ACCESS_FAILURE"
	CodeDBInsertFailure = "DB_INSERT_FAILURE"
	CodeDBUpdateFailure = "DB_UPDATE_FAILURE"
	CodeDBRemoveFailure = "DB_REMOVE_FAILURE"

	CodeDuplicateEntry      = "DUPLICATE_ENTRY"
	CodeForeignKeyViolation = "FOREIGN_KEY_VIOLATION"

	CodeProductNotFound = "PRODUCT_NOT_FOUND"
	CodeUserNotFound    = "USER_NOT_FOUND"
	CodeGroupNotFound   = "GROUP_NOT_FOUND"
	CodeInviteNotFound  = "INVITE_NOT_FOUND"

	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeSessionFailure     = "SESSION_FAILURE"
	CodeSessionSaveFailure = "SESSION_SAVE_FAILURE"

	CodeInviteSelf              = "INVITE_SELF"
	CodeAlreadyInvited          = "ALREADY_INVITED"
	CodeAlreadyMember           = "ALREADY_MEMBER"
	CodeNotAMember              = "NOT_A_MEMBER"
	CodeInsufficientPermissions = "INSUFFICIENT_PERMISSIONS"
	CodeUsernameTaken           = "USERNAME_TAKEN"

	CodeRequired         = "REQUIRED"
	CodeMaxLength        = "MAX_LENGTH"
	CodeInvalidURL       = "INVALID_URL"
	CodeAlphaSpace       = "ALPHA_SPACE"
	CodeInvalidDate      = "INVALID_DATE"
	CodeValidationFailed = "VALIDATION_FAILED"
)

type ErrorItem struct {
	Code   string         `json:"code"`
	Detail string         `json:"detail"`
	Meta   map[string]any `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Errors []ErrorItem `json:"errors"`
}

func NewError(code, detail string) ErrorItem {
	return ErrorItem{Code: code, Detail: detail}
}

func NewErrorWithMeta(code, detail string, meta map[string]any) ErrorItem {
	return ErrorItem{Code: code, Detail: detail, Meta: meta}
}

func writeJSON(w http.ResponseWriter, status int, errors ...ErrorItem) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Errors: errors})
}

func ServerError(w http.ResponseWriter, errors ...ErrorItem) {
	writeJSON(w, http.StatusInternalServerError, errors...)
}

func BadRequest(w http.ResponseWriter, errors ...ErrorItem) {
	writeJSON(w, http.StatusBadRequest, errors...)
}

func ValidationErrors(w http.ResponseWriter, errors ...ErrorItem) {
	writeJSON(w, http.StatusUnprocessableEntity, errors...)
}

func NotFound(w http.ResponseWriter, errors ...ErrorItem) {
	writeJSON(w, http.StatusNotFound, errors...)
}

func Unauthorized(w http.ResponseWriter, errors ...ErrorItem) {
	writeJSON(w, http.StatusUnauthorized, errors...)
}

func Forbidden(w http.ResponseWriter, errors ...ErrorItem) {
	writeJSON(w, http.StatusForbidden, errors...)
}

func Conflict(w http.ResponseWriter, errors ...ErrorItem) {
	writeJSON(w, http.StatusConflict, errors...)
}

// ClassifyDBError maps well-known Postgres constraint violations to specific
// error items. It reports false when the error is not one it recognizes, in
// which case the caller should fall back to a generic response.
func ClassifyDBError(err error) (ErrorItem, bool) {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return ErrorItem{}, false
	}

	switch pgErr.Code {
	case "23505": // unique_violation
		return NewErrorWithMeta(CodeDuplicateEntry, "duplicate entry", map[string]any{
			"constraint": pgErr.ConstraintName,
		}), true
	case "23503": // foreign_key_violation
		return NewErrorWithMeta(CodeForeignKeyViolation, "referenced record does not exist", map[string]any{
			"constraint": pgErr.ConstraintName,
		}), true
	}

	return ErrorItem{}, false
}

// WriteDBError writes a classified constraint violation with its proper status
// (409 for duplicates, 422 for foreign key violations), or the fallback item
// as a 500 when the error is not a recognized constraint violation.
func WriteDBError(w http.ResponseWriter, err error, fallback ErrorItem) {
	item, ok := ClassifyDBError(err)
	if !ok {
		ServerError(w, fallback)
		return
	}

	switch item.Code {
	case CodeDuplicateEntry:
		Conflict(w, item)
	case CodeForeignKeyViolation:
		ValidationErrors(w, item)
	default:
		ServerError(w, item)
	}
}
