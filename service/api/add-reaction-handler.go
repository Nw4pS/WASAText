package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// payload atteso
type addReactionRequest struct {
	Emoji string `json:"emoji"`
}

func (rt *_router) addReaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// verifica token auth
	authHeader := r.Header.Get("Authorization")
	requesterID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if requesterID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// estrazione id dal path
	msgIDStr := ps.ByName("id")
	messageID, err := strconv.Atoi(msgIDStr)
	if err != nil {
		ctx.Logger.WithError(err).Warning("ID messaggio non valido")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// decodifica json
	var req addReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la decodifica del JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validazione emoji
	emojiStr := strings.TrimSpace(req.Emoji)
	emojiLen := utf8.RuneCountInString(emojiStr)

	if emojiLen < 1 || emojiLen > 8 {
		ctx.Logger.Warning("Lunghezza emoji non conforme (richiesta tra 1 e 8)")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// chiamata al db
	err = rt.db.AddReaction(requesterID, messageID, emojiStr)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante l'aggiunta della reazione")

		errStr := err.Error()
		if strings.Contains(errStr, "operazione negata") {
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
