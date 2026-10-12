package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// oggetto in risposta
type conversationResponse struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Picture            string `json:"picture"`
	LastMessageDate    string `json:"lastMessageDate"`
	LastMessagePreview string `json:"lastMessagePreview"`
	IsGroup            bool   `json:"isGroup"`
}

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// estrazione del token auth
	authHeader := r.Header.Get("Authorization")
	myID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if myID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// chiamata al database
	convList, err := rt.db.GetMyConversations(myID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante il recupero delle conversazioni")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// se la lista è vuota, restituisco un array JSON vuoto [] anziché null
	resp := make([]conversationResponse, 0, len(convList))
	for _, c := range convList {
		resp = append(resp, conversationResponse{
			ID:                 c.ID,
			Name:               c.Name,
			Picture:            c.Picture,
			LastMessageDate:    c.LastMessageDate,
			LastMessagePreview: c.LastMessagePreview,
			IsGroup:            c.IsGroup,
		})
	}

	// risposta 204 successo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
