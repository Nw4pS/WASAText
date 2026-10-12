package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type propicRequest struct {
	Picture string `json:"picture"`
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var propic propicRequest

	if err := json.NewDecoder(r.Body).Decode(&propic); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if propic.Picture == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	authHeader := r.Header.Get("Authorization")
	id := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if id == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err := rt.db.SetPropic(id, propic.Picture)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante l'aggiornamento della foto profilo nel database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
