package inviteHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sl-api/api/routes/common/errors"
	"sl-api/api/services/group"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *API) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	groupID, err := uuid.Parse(chi.URLParam(r, "groupId"))
	if err != nil {
		errorsCommon.BadRequest(w, errorsCommon.RespInvalidUUID)
		return
	}

	input := groupService.AcceptInviteInput{
		GroupID: groupID,
	}

	output, err := api.acceptInviteService.Execute(r.Context(), input)

	if err != nil {
		fmt.Println("Error executing accept invite service:", err)

		errorsCommon.ServerError(w, errorsCommon.RespGenericFailure)
		return

		// TODO: Create the correct errors
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		fmt.Println("Error encoding accept invite output:", err)
		errorsCommon.ServerError(w, errorsCommon.RespJSONEncodeFailure)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
