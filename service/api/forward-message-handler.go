package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// body della richiesta
type forwardMessageRequest struct {
	MessageID int `json:"messageId"`
}

// risposta
type forwardMessageResponse struct {
	ID int `json:"id"`
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// verifica esistenza token auth
	authHeader := r.Header.Get("Authorization")
	myID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if myID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// estrazione id conversazione dal path
	destConvIDStr := ps.ByName("id")
	destConvID, err := strconv.Atoi(destConvIDStr)
	if err != nil {
		ctx.Logger.WithError(err).Warning("ID conversazione non valido")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// decodifica corpo json
	var req forwardMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la decodifica del JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// chiamata funzione database
	newMsgID, err := rt.db.ForwardMessage(myID, destConvID, req.MessageID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante l'inoltro del messaggio")

		errStr := err.Error()
		if strings.Contains(errStr, "operazione negata") {
			w.WriteHeader(http.StatusForbidden)
			return
		} else if strings.Contains(errStr, "messaggio originale non trovato") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// preparo e invio della risposta
	resp := forwardMessageResponse{
		ID: newMsgID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la codifica della risposta JSON")
	}
}
