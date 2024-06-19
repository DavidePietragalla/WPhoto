package api

import (
	"WasaPhotoDavidePietragalla/service/api/reqcontext"
	"WasaPhotoDavidePietragalla/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// This function retrieves all the posts of the people that the user is following
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// imposto il tipo di contenuto della risposta http in json
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

	// Ottengo un elenco di utenti che l'utente sta seguendo dal database.
	followers, err := rt.db.GetFollowing(requestingUserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Itera sugli utenti seguiti e recupera le loro foto. Aggiungendo le foto a un elenco.
	var posts []database.Post
	for _, follower := range followers {

		followerPost, err := rt.db.GetPostsList(requestingUserId, follower.UserId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for i, post := range followerPost {
			if i >= database.PostsPerUserHome {
				break
			}
			posts = append(posts, post)
		}

	}

	// Imposta lo stato della risposta HTTP come 200 OK. Codifica l'elenco di foto in formato JSON e lo invia come corpo della risposta HTTP.
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(posts)
}
