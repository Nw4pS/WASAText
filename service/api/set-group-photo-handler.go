package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type setGroupPhotoRequest struct {
	Picture string `json:"picture"`
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// verifica della presenza id auth
	authHeader := r.Header.Get("Authorization")
	myID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if myID == "" {
		w.WriteHeader(http.StatusUnauthorized) // 401 Unauthorized
		return
	}

	// parsing e validazione input
	groupIDStr := ps.ByName("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		ctx.Logger.WithError(err).Warning("ID gruppo non valido per setGroupPhoto")
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		return
	}

	// decodifica json
	var req setGroupPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la decodifica del corpo JSON")
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		return
	}

	// la stringa picture non deve essere vuota
	trimmedPicture := strings.TrimSpace(req.Picture)
	if trimmedPicture == "" {
		ctx.Logger.Warning("Stringa immagine vuota")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validazione del formato base64
	if _, err := base64.StdEncoding.DecodeString(trimmedPicture); err != nil {
		ctx.Logger.WithError(err).Warning("Formato Base64 non valido")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// chiamata al database
	err = rt.db.SetGroupPhoto(myID, groupID, trimmedPicture)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore del database durante l'aggiornamento della foto del gruppo")

		if strings.Contains(err.Error(), "operazione negata") {
			w.WriteHeader(http.StatusForbidden) // 403 Forbidden: l'utente non fa parte del gruppo
			return
		}

		if strings.Contains(err.Error(), "gruppo non trovato") {
			w.WriteHeader(http.StatusNotFound) // 404 Not Found
			return
		}

		w.WriteHeader(http.StatusInternalServerError) // 500 Errore generico
		return
	}

	// 204 No Content come definito in api.yaml
	w.WriteHeader(http.StatusNoContent)
}
