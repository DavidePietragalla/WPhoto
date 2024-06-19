package api

import (
	"WasaPhotoDavidePietragalla/service/api/reqcontext"
	"WasaPhotoDavidePietragalla/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// Funzione che ritrova tutte le info necessarie del profilo
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	var followers []database.User
	var following []database.User
	var posts []database.Post

	// Controllo se l'utente esiste
	userExists, err := rt.db.CheckUser(nickname)
	if err != nil {
		ctx.Logger.WithError(err).Error("getProfile/db.CheckUser: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !userExists {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Recupero l'id dell'utente richiesto
	requestedUserId, err := rt.db.GetId(nickname)
	if err != nil {
		ctx.Logger.WithError(err).Error("getProfile/db.GetId: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Controlla se l'utente che effettua la richiesta è bannato dall'utente richiesto, utilizzando una funzione del database.
	// Se l'utente che effettua la richiesta è bannato, la funzione restituisce un codice di stato HTTP 403 (Forbidden) e termina.
	userBanned, err := rt.db.BannedUserCheck(requestingUserId, requestedUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("getProfile/db.BannedUserCheck/userBanned: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userBanned {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Recupero la lista dei followers dell'utente richiesto
	followers, err = rt.db.GetFollowers(requestedUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.GetFollowers: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Recupero la lista dei utenti seguiti dall'utente richiesto
	following, err = rt.db.GetFollowing(requestedUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.GetFollowing: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Recupera la lista delle foto dell'utente richiesto dal database.
	posts, err = rt.db.GetPostsList(requestingUserId, requestedUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("getProfile/db.GetPostsList: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Imposta il codice di stato della risposta HTTP come 200 (OK) e invia un oggetto
	// JSON che rappresenta il profilo completo dell'utente richiesto come corpo della risposta.
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Profile{
		UserId:    requestedUserId,
		Nickname:  nickname,
		Followers: followers,
		Following: following,
		Posts:     posts,
	})

}

// Funzione per cambiare nickname(gestisce le richieste http put)
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Estraggo il nuovo nickname dal corpo della richiesta e decodifica del JSON
	var nick Nickname
	err := json.NewDecoder(r.Body).Decode(&nick)
	// Se c'è un errore nella decodifica del JSON, si risponde con un codice di stato HTTP 400 (Bad Request).
	if err != nil {
		ctx.Logger.WithError(err).Error("update-nickname: error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Modifico il nickname dell'utente nel database
	err = rt.db.ChangeNickname(
		requestingUserId,
		nick.Nickname)
	if err != nil {
		ctx.Logger.WithError(err).Error("update-nickname: error executing update query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta l'intestazione della risposta per indicare che il tipo di contenuto sarà JSON.
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

	// Estraggo il parametro di query "requestedId" dall'URL della richiesta.
	requestedId := r.URL.Query().Get("requestedId")
	requestedUserId, _ := strconv.ParseInt(requestedId, 10, 64)

	// ottengo il nickname dell'utente nel database
	nick, err := rt.db.GetNickname(requestedUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("get-nickname: error executing select query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 200 http status
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Nickname{
		Nickname: nick,
	})

}
