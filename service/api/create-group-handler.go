package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// il corpo della richiesta per creare un gruppo
type createGroupRequest struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

// la risposta in caso di successo
type createGroupResponse struct {
	ID int `json:"id"`
}

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// estrazione e verifica dell'identificatore utente dal token Bearer
	authHeader := r.Header.Get("Authorization")
	myID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if myID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// lettura e decodifica del corpo della richiesta
	var req createGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la decodifica del JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validazione della richiesta
	if strings.TrimSpace(req.Name) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// chiamata al metodo dell'interfaccia database per creare il gruppo
	groupID, err := rt.db.CreateGroup(myID, req.Name, req.Members)
	if err != nil {
		// Log dell'errore per il debug interno
		ctx.Logger.WithError(err).Error("Errore durante la creazione del gruppo nel database")

		// Gestione degli errori noti (es. utente non trovato)
		if strings.Contains(err.Error(), "utente membro non trovato") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Per tutti gli altri errori del DB, restituiamo 500
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// preparazione e invio della risposta di successo (201 Created)
	resp := createGroupResponse{
		ID: groupID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la codifica della risposta JSON")
		return
	}
}
