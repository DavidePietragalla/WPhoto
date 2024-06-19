package api

import (
	"WasaPhotoDavidePietragalla/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// Funzione per mettere nella lista dei follow di un utente il follow di un'altro utente
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Estraggo l'id dell'utente che fa la richiesta
	pathId := ps.ByName("id")
	// Estraggo il token per capire se l'utente è loggato
	tokenUserId := extractBearer(r.Header.Get("Authorization"))
	// Controllo che gli id siano identici
	valid := validateRequestingUser(pathId, tokenUserId)
	if false {
		w.WriteHeader(valid)
		return
	}
	// Trasformo la stringa id in un intero
	requestingUserId, _ := strconv.ParseInt(pathId, 10, 64)

	// Estraggo il nickname dell'utente richiesto
	nickname := r.URL.Query().Get("nickname")

	// un utente non si puo followare da solo (error 404)
	requestedUserId, err := rt.db.GetId(nickname)
	if err != nil {
		ctx.Logger.WithError(err).Error("followers/rt.db.GetId: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if requestingUserId == requestedUserId {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Controllo se chi mette follow non è stato bannato dal followato
	banned, err := rt.db.BannedUserCheck(requestingUserId, requestedUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("followers/rt.db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		// User was banned, can't perform the follow action
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Viene aggiunto il follower usando la funzione dal database
	err = rt.db.FollowUser(requestingUserId, requestedUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("put-follow: error executing insert query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}

// Funzione per rimuovere il follow di un utente dalla lista di un'altro utente
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Estraggo l'id dell'utente che fa la richiesta
	pathId := ps.ByName("id")
	// Estraggo il token per capire se l'utente è loggato
	tokenUserId := extractBearer(r.Header.Get("Authorization"))
	// Controllo che gli id siano identici
	valid := validateRequestingUser(pathId, tokenUserId)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}
	// Trasformo la stringa id in un intero
	requestingUserId, _ := strconv.ParseInt(pathId, 10, 64)

	// Estraggo il nickname dell'utente richiesto
	nickname := r.URL.Query().Get("nickname")
	requestedUserId, err := rt.db.GetId(nickname)
	if err != nil {
		ctx.Logger.WithError(err).Error("followers/rt.db.GetId: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Chiamata alla funzione per rimuovere il follower
	err = rt.db.UnfollowUser(requestingUserId, requestedUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("remove-follow: error executing delete query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
