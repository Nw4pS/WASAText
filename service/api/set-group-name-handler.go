package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

type setGroupNameRequest struct {
	Name string `json:"name"`
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// estraggo il token di autenticazione
	authHeader := r.Header.Get("Authorization")
	requesterID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if requesterID == "" {
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}

	// parsing del parametro id del gruppo dal path
	groupID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil || groupID <= 0 {
		ctx.Logger.WithError(err).Warning("ID gruppo non valido per setGroupName")
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	// decodifica json
	var req setGroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Warning("Body JSON non valido per setGroupName")
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	trimmedName := strings.TrimSpace(req.Name)
	if trimmedName == "" {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	// chiamata al db
	err = rt.db.SetGroupName(requesterID, groupID, trimmedName)
	if err != nil {
		if errors.Is(err, database.ErrNotGroupMember) {
			w.WriteHeader(http.StatusForbidden) // 403
			return
		}
		if errors.Is(err, database.ErrGroupNotFound) {
			w.WriteHeader(http.StatusNotFound) // 404
			return
		}
		if errors.Is(err, database.ErrInvalidName) {
			w.WriteHeader(http.StatusBadRequest) // 400
			return
		}

		ctx.Logger.WithError(err).Error("Errore interno durante il cambio nome del gruppo")
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}

	// risposta 204
	w.WriteHeader(http.StatusNoContent)
}
