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
type addToGroupRequest struct {
	Usernames []string `json:"usernames"`
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// estrazione e verifica dell'identificatore utente
	authHeader := r.Header.Get("Authorization")
	myID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if myID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// estrazione e validazione dell'id del gruppo dal path
	groupIDStr := ps.ByName("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		ctx.Logger.WithError(err).Warning("ID gruppo non valido")
		w.WriteHeader(http.StatusBadRequest) // 400 se l'id non è un numero intero
		return
	}

	// lettura e decodifica del corpo della richiesta JSON
	var req addToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la decodifica del JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// l'array degli username non deve essere vuoto
	if len(req.Usernames) == 0 {
		ctx.Logger.Warning("Array di username vuoto")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// controllo che non ci siano stringhe vuote nell'array
	for _, un := range req.Usernames {
		if strings.TrimSpace(un) == "" {
			ctx.Logger.Warning("Username vuoto fornito nell'array")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// chiamata all'interfaccia del database
	err = rt.db.AddToGroup(myID, groupID, req.Usernames)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante l'aggiunta degli utenti al gruppo")

		// gestione degli errori noti provenienti dal layer database
		errStr := err.Error()
		if strings.Contains(errStr, "operazione negata") {
			// g'utente richiedente non è membro del gruppo
			w.WriteHeader(http.StatusForbidden) // 403 Forbidden è più appropriato qui, o 401
			return
		} else if strings.Contains(errStr, "utente membro non trovato") {
			// uno degli utenti da aggiungere non esiste
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if strings.Contains(errStr, "fa già parte di questo gruppo") {
			// utenti già presenti
			w.WriteHeader(http.StatusBadRequest) // O 409 Conflict
			return
		}

		// per tutti gli altri errori, si restituisce 500
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// successo
	w.WriteHeader(http.StatusNoContent)
}
