package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// payload aspettato
type sendMessageRequest struct {
	Content string `json:"content"`
	IsImage bool   `json:"isImage"`
	ReplyTo *int   `json:"replyTo,omitempty"`
}

// risposta in caso di successo
type sendMessageResponse struct {
	ID int `json:"id"`
}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// verifico che il token auth sia presente
	authHeader := r.Header.Get("Authorization")
	senderID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if senderID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// estraggo l'id della conversazione dal path
	convIDStr := ps.ByName("id")
	convID, err := strconv.Atoi(convIDStr)
	if err != nil {
		ctx.Logger.WithError(err).Warning("ID conversazione non valido")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// decodifico il json
	var req sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Errore decodifica JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// controllo che il contenuto non sia vuoto
	if strings.TrimSpace(req.Content) == "" {
		ctx.Logger.Warning("Contenuto del messaggio vuoto")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// siccome sono tanti parametri, conviene creare una struct perché in caso di future modifiche è più facile da modificare
	dbReq := database.SendMessageRequest{
		ConversationID: convID,
		SenderID:       senderID,
		Content:        req.Content,
		IsImage:        req.IsImage,
		ReplyTo:        req.ReplyTo, // sarà un puntatore, in caso non punti a nulla sarà null
	}

	// chiamata all'interfaccia db
	messageID, err := rt.db.SendMessage(dbReq)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante l'invio del messaggio")

		if strings.Contains(err.Error(), "operazione negata") {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// prep risposta
	resp := sendMessageResponse{
		ID: messageID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		ctx.Logger.WithError(err).Error("Errore codifica risposta JSON")
	}
}
