package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// estraggo e verifico l'identificatore utente
	authHeader := r.Header.Get("Authorization")
	myID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if myID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// estraggo e valido l'd del gruppo dal path
	groupIDStr := ps.ByName("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		ctx.Logger.WithError(err).Warning("ID gruppo non valido per getGroupMembers")
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		return
	}

	// chiamo l'interfaccia del database
	members, err := rt.db.GetGroupMembers(myID, groupID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante il recupero dei membri del gruppo")

		// utente non membro del gruppo
		if strings.Contains(err.Error(), "operazione negata") {
			w.WriteHeader(http.StatusForbidden) // 403 Forbidden
			return
		}

		// errore generico del server
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// imposto gli headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// scrivo il body della risposta (la lista utenti)
	if err := json.NewEncoder(w).Encode(members); err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la codifica della risposta JSON")
		return
	}
}
