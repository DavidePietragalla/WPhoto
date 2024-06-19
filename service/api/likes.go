package api

import (
	"WasaPhotoDavidePietragalla/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// Funzione permette a un utente di mettere like ad una foto
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Estraggo l'id del post richiesto
	stringPostId := r.URL.Query().Get("postId")
	postId, _ := strconv.ParseInt(stringPostId, 10, 64)

	// Chiama una funzione del database per aggiungere il "like" alla foto specificata.
	err := rt.db.LikePost(postId, requestingUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("likePost: error executing insert query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}

// Funzione per rimuovere il like da una foto
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Estraggo l'id del post richiesto
	stringPostId := r.URL.Query().Get("postId")
	postId, _ := strconv.ParseInt(stringPostId, 10, 64)

	// Chiamata della funzione per rimuovere il "like".
	err := rt.db.UnlikePost(postId, requestingUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("unlikePost/db.UnlikePost: error executing insert query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
