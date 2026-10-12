package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// searchUsers gestisce l'endpoint GET /users
func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// validazione del token Bearer per l'autenticazione
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// estrazione del parametro di ricerca opzionale
	searchQuery := r.URL.Query().Get("username")

	// chiamata al livello database
	users, err := rt.db.SearchUsers(searchQuery)
	if err != nil {
		ctx.Logger.WithError(err).Error("errore del database durante la ricerca utenti")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// configurazione e invio della risposta HTTP 200 OK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		ctx.Logger.WithError(err).Error("errore durante la serializzazione JSON della risposta")
	}
}
