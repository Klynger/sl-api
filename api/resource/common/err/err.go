package err

import (
	"net/http"
)

var (
	RespDBDataInsertFailure  = []byte(`{"error": "db data insert failure"}`)
	RespDBDataAccessFailure  = []byte(`{"error": "db data access failure"}`)
	RespDBDataUpdateFailure  = []byte(`{"error": "db data update failure"}`)
	RespDBDdataRemoveFailure = []byte(`{"error": "db data remove failure"}`)

	RespJSONEncodeFailure = []byte(`{"error": "json encode failure"}`)
	RespJSONDecodeFailure = []byte(`{"error": "json decode failure"}`)

	RespInvalidURLParamID = []byte(`{"error": "invalid url param-id"}`)

	RespSessionAccessFailure = []byte(`{"error": "session access failure"}`)

	RespGenericFailure = []byte(`{"error": "something went wrong"}`)

	RespUnauthorized = []byte(`{"error": "unauthorized"}`)

	RespInvalidUUID = []byte(`{"code": "INVALID_UUID", "error": "invalid uuid"}`)

	// TODO: Create specific errors in the service layer and return them here instead of generic server error
	RespGroupNotFound = []byte(`{"code": "GROUP_NOT_FOUND", "error": "group not found"}`)
)

func ServerError(w http.ResponseWriter, reps []byte) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(reps)
}

func BadRequest(w http.ResponseWriter, reps []byte) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write(reps)
}

func ValidationErrors(w http.ResponseWriter, reps []byte) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write(reps)
}

type Error struct {
	Error string `json:"error"`
}

type Errors struct {
	Errors []string `json:"errors"`
}
