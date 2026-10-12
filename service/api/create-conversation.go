package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// createConvRequest rappresenta il body JSON in ingresso
type createConvRequest struct {
	Username string `json:"username"`
}

// createConvResponse rappresenta il JSON di risposta con status 201
type createConvResponse struct {
	ID int `json:"id"`
}

func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// estrazione e verifica dell'identificatore utente dal token Bearer
	authHeader := r.Header.Get("Authorization")
	myID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if myID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// lettura e validazione del body della richiesta
	var req createConvRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// controllo che l'username non sia vuoto
	if req.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// chiamata al livello database per creare o recuperare la conversazione
	convID, err := rt.db.CreateConversation(myID, req.Username)
	if err != nil {
		// loggo l'errore per il debug interno usando il logger di contesto
		ctx.Logger.WithError(err).Error("Errore durante la creazione/recupero della conversazione nel database")
		// restituisco 500 Internal Server Error al client
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// risposta di successo: 201 Created con l'ID della conversazione
	resp := createConvResponse{
		ID: convID,
	}

	// imposto l'header Content-Type per indicare che stiamo restituendo un JSON
	w.Header().Set("Content-Type", "application/json")
	// scrivo lo status code richiesto dalle specifiche
	w.WriteHeader(http.StatusCreated)

	// serializzo la struct resp nel body della risposta HTTP
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		ctx.Logger.WithError(err).Error("Errore durante l'encoding della risposta JSON")
		return
	}
}
