package api

import (
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// verifico che il token sia stato inserito

	authHeader := r.Header.Get("Authorization")
	requesterID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if requesterID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// estraggo id messaggio dal path
	msgIDStr := ps.ByName("id")
	messageID, err := strconv.Atoi(msgIDStr)
	if err != nil {
		ctx.Logger.WithError(err).Warning("ID messaggio non valido")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// chiamo funzione db
	err = rt.db.DeleteMessage(requesterID, messageID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante l'eliminazione del messaggio")

		errStr := err.Error()
		if strings.Contains(errStr, "operazione negata") {
			// se l'utente richiedente l'eliminazione non è il mittente
			w.WriteHeader(http.StatusForbidden)
			return
		} else if strings.Contains(errStr, "messaggio non trovato") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
