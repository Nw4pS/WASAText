package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteReaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authHeader := r.Header.Get("Authorization")
	requesterID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if requesterID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	messageID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil || messageID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteReaction(requesterID, messageID)
	if err != nil {
		if errors.Is(err, database.ErrReactionNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("errore rimozione reazione")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
