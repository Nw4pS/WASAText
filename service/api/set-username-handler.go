package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// il body JSON atteso da api.yaml
type setMyUserNameRequest struct {
	Name string `json:"name"`
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// estrazione e verifica dell'identificatore utente dal token Bearer
	authHeader := r.Header.Get("Authorization")
	myID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if myID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// lettura e decodifica del corpo della richiesta JSON
	var req setMyUserNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la decodifica del JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validazione dell'input: minLength 3, maxLength 16 come da api.yaml
	nameStr := strings.TrimSpace(req.Name)
	nameLen := utf8.RuneCountInString(nameStr)
	if nameLen < 3 || nameLen > 16 {
		ctx.Logger.Warning("Lunghezza username non conforme (richiesta tra 3 e 16 caratteri)")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// chiamata alla funzione del database
	err := rt.db.SetMyUsername(myID, nameStr)
	if err != nil {
		// controllo username già occupato
		if errors.Is(err, database.ErrUsernameTaken) {
			ctx.Logger.Warning("Tentativo di impostare un username già in uso")
			w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
			return
		}

		// controllo utente del token non trovato
		if errors.Is(err, database.ErrUserNotFound) {
			ctx.Logger.Warning("Utente richiedente non trovato nel database")
			w.WriteHeader(http.StatusUnauthorized) // 401 Unauthorized
			return
		}

		// errore generico
		ctx.Logger.WithError(err).Error("Errore durante l'aggiornamento dell'username nel database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// successo 204
	w.WriteHeader(http.StatusNoContent)
}
