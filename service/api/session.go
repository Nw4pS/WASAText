package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// LoginRequest mappa il payload in ingresso definito nell'API
type LoginRequest struct {
	Name string `json:"name"`
}

// LoginResponse mappa la risposta contenente l'identificatore
type LoginResponse struct {
	Identifier string `json:"identifier"`
}

// doLogin è l'handler dell'endpoint POST /session.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req LoginRequest

	// decodifica il body della richiesta, risponde con 400 se il JSON è malformato
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// l'username deve essere compreso tra i 3 e i 16 caratteri
	name := strings.TrimSpace(req.Name)
	if len(name) < 3 || len(name) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// delega al livello database la logica di business.
	identifier, err := rt.db.DoLogin(name)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore del database durante il doLogin")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// prepara l'oggetto JSON di risposta con l'ID (generato o recuperato)
	resp := LoginResponse{
		Identifier: identifier,
	}

	// imposta gli header e lo status code obbligatorio 201 Created
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// invia la risposta codificata al client.
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		ctx.Logger.WithError(err).Error("Errore durante l'encoding della risposta")
	}
}
