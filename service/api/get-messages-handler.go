package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Estrazione del Bearer Token dell'utente richiedente
	authHeader := r.Header.Get("Authorization")
	requesterID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if requesterID == "" {
		ctx.Logger.Warning("Token Bearer assente o vuoto")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 2. Estrazione e validazione dell'ID della conversazione dal path URL
	convIDStr := ps.ByName("id")
	convID, err := strconv.Atoi(convIDStr)
	if err != nil {
		ctx.Logger.WithError(err).Warning("ID conversazione non valido")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 3. Chiamata al database
	messages, err := rt.db.GetConversation(requesterID, convID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante il recupero dei messaggi della conversazione")

		// Se l'utente non fa parte della chat, restituiamo 403 Forbidden
		if strings.Contains(err.Error(), "operazione negata") {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 4. Risposta di successo con status 200 OK e array JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la codifica JSON dei messaggi")
	}
}
