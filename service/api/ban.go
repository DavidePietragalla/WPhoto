package api

import (
	"WasaPhotoDavidePietragalla/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// Funzione che gestisce l'aggiunta di un utente alla lista degli utenti bannati di un altro utente.
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Controlla se l'utente sta cercando di bannare se stesso.se si(400:"Bad Request")
	requestedUserId, err := rt.db.GetId(nickname)
	if err != nil {
		ctx.Logger.WithError(err).Error("ban/rt.db.GetId: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if requestingUserId == requestedUserId {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Chiama una funzione del database per aggiungere l'utente specificato alla lista degli utenti bannati
	err = rt.db.BanUser(requestingUserId, requestedUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("put-ban/db.BanUser: error executing insert query")

		// C'è stato un errore interno,restituisco(error:500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Bannare implica anche rimuovere il follow (se esiste)
	err = rt.db.UnfollowUser(requestingUserId, requestedUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("put-ban/db.UnfollowUser1: error executing insert query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}

// Funzione che rimuove un user dalla lista dei banned di un altro user
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Recupero l'id dell'utente da sbannare
	requestedUserId, err := rt.db.GetId(nickname)
	if err != nil {
		ctx.Logger.WithError(err).Error("unban/rt.db.GetId: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Chiamo la funzione UnbanUser dal database per rimuovere l'utente dalla lista dei banned.
	err = rt.db.UnbanUser(requestingUserId, requestedUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("unban: error executing delete query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
