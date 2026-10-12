package api

import (
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// estrazione e verifica dell'id utente
	authHeader := r.Header.Get("Authorization")
	myID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if myID == "" {
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}

	// estrazione e validazione dell'id preso dal path
	groupIDStr := ps.ByName("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		ctx.Logger.WithError(err).Warning("ID gruppo non valido per leaveGroup")
		w.WriteHeader(http.StatusBadRequest) // 400 se l'id non è un numero intero
		return
	}

	// chiamata all'interfaccia database
	err = rt.db.LeaveGroup(myID, groupID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante l'abbandono del gruppo")

		if strings.Contains(err.Error(), "non fa parte di questo gruppo") {
			// l'utente ha tentato di uscire da un gruppo in cui non è presente
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// per tutti gli altri problemi del database
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// '204': description: Uscito dal gruppo con successo
	w.WriteHeader(http.StatusNoContent)
}
