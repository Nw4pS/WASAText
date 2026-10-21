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
	// 1. Verifica della presenza del token di autorizzazione (Bearer Auth di WASAText)
	authHeader := r.Header.Get("Authorization")
	myID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if myID == "" {
		w.WriteHeader(http.StatusUnauthorized) // 401 Unauthorized
		return
	}

	// 2. Parsing e validazione dell'ID del gruppo
	groupIDStr := ps.ByName("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil || groupID <= 0 {
		ctx.Logger.WithError(err).Warning("ID gruppo non valido per setGroupPhoto")
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		return
	}

	// 3. Decodifica del corpo JSON
	var req setGroupPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la decodifica del corpo JSON")
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		return
	}

	// 4. Verifica che la stringa immagine non sia vuota
	trimmedPicture := strings.TrimSpace(req.Picture)
	if trimmedPicture == "" {
		ctx.Logger.Warning("Stringa immagine vuota")
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		return
	}

	// 5. Validazione del formato Base64:
	// Se la stringa proviene dal browser (FileReader.readAsDataURL), include l'intestazione
	// MIME (es. "data:image/png;base64,"). Estraiamo la parte successiva alla virgola
	// per la verifica formale, mantenendo l'intera stringa per il salvataggio su DB.
	b64Payload := trimmedPicture
	if commaIdx := strings.Index(trimmedPicture, ","); commaIdx != -1 {
		b64Payload = trimmedPicture[commaIdx+1:]
	}

	if _, err := base64.StdEncoding.DecodeString(b64Payload); err != nil {
		ctx.Logger.WithError(err).Warning("Formato Base64 non valido")
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		return
	}

	// 6. Chiamata al layer database
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

	// 7. Risposta 204 No Content conforme ad api.yaml
	w.WriteHeader(http.StatusNoContent)
}
