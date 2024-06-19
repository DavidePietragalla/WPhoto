package api

import (
	"WasaPhotoDavidePietragalla/service/api/reqcontext"
	"WasaPhotoDavidePietragalla/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Funzione che gestisce l'upload di una foto
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// recupero la data attuale
	date := time.Now().UTC()

	// Legge il body della richiesta e verifica se ci sono errori durante la lettura.
	var err error
	var urlImage Url
	urlImage.Data, err = io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("uploadPost: error reading body content")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Chiama una funzione del database per creare un record per la foto e ottenere un ID univoco per essa
	// Se si verifica un errore, risponde con un codice di stato HTTP 500 (Internal Server Error)
	postIdInt, err := rt.db.CreatePost(requestingUserId, date, urlImage.Data)
	if err != nil {
		ctx.Logger.WithError(err).Error("uploadPost: error executing db function call")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Invia una risposta con stato "Created" e un oggetto JSON che rappresenta la foto appena caricata.
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(database.Post{
		PostId:   postIdInt,
		UserId:   requestingUserId,
		Date:     date,
		Url:      urlImage.Data,
		Comments: []database.Comment{},
		Likes:    []database.User{},
	})

}

// Funzione che restituisce la foto richiesta
// Vengono estratti user_id,post_id si crea il percorso
// e infine viene servito il file con il metodo http
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Recupero il post
	post, err := rt.db.GetPost(requestingUserId, postId)
	if err != nil {
		ctx.Logger.WithError(err).Error("getPost: error executing GetPost")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Invia una risposta con stato "Created" e un oggetto JSON che rappresenta la foto appena caricata.
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(post)
}

// Function that deletes a post (this includes comments and likes)
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Estraggo l'id del post da eliminare
	stringPostId := r.URL.Query().Get("postId")
	postId, _ := strconv.ParseInt(stringPostId, 10, 64)

	// Chiama una funzione per rimuovere la foto dal database.
	err := rt.db.RemovePost(requestingUserId, postId)
	if err != nil {
		ctx.Logger.WithError(err).Error("post-delete/RemovePost: error coming from database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
