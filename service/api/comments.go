package api

import (
	"WasaPhotoDavidePietragalla/service/api/reqcontext"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// funzione che gestisce l'aggiunta di un commento a una foto e invia una risposta contenente l'ID univoco del commento creato.
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

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

	// Recupero l'id dell'autore del post
	authorId, _ := rt.db.GetPostAuthor(postId)

	// Controlla se l'utente che effettua la richiesta è stato bannato dal proprietario della foto.
	banned, err := rt.db.BannedUserCheck(requestingUserId, authorId)
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPost/db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		// l'utente bannato non puo commentare
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Copio il body content(comment mandato dal'user) nel comment(Struct)
	var comment Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("commentPost/Decode: failed to decode request body json")
		return
	}

	// Controllo la lunghezza del comment(<=30)
	if len(comment.Comment) > 30 {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("commentPost: comment longer than 30 characters")
		return
	}

	// Chiama una funzione del database per creare il commento.
	commentId, err := rt.db.CommentPost(postId, requestingUserId, comment.Comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("commentPost/db.CommentPost: failed to execute query for insertion")
		return
	}

	// Imposta lo stato della risposta su "Created".
	w.WriteHeader(http.StatusCreated)

	// Codifica l'ID univoco del commento creato in formato JSON e lo invia come corpo della risposta.
	err = json.NewEncoder(w).Encode(CommentId{CommentId: commentId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("commentPost/Encode: failed convert post_id to int64")
		return
	}
}

// Funzione che rimuove un commento da una foto
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Estraggo l'id del commento richiesto
	stringCommentId := r.URL.Query().Get("commentId")
	commentId, _ := strconv.ParseInt(stringCommentId, 10, 64)

	// Controllo per vedere se l'utente che effettua la richiesta è l'autore del post.
	authorId, _ := rt.db.GetCommentAuthor(commentId)
	if requestingUserId == authorId {
		// Chiamo la funzione dal db per rimuovere il commento
		err := rt.db.UncommentPost(requestingUserId, commentId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("uncommentPost: failed to execute query for removing")
			return
		}
		// Rispondo con codice 204(la richiesta è andata a buon fine)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// L'utente non è l'autore
	w.WriteHeader(http.StatusForbidden)
}
