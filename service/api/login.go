package api

import (
	"WasaPhotoDavidePietragalla/service/api/reqcontext"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Funzione che gestice la sessione degli utenti
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta l'intestazione della risposta per indicare che il tipo di contenuto sarà JSON.
	w.Header().Set("Content-Type", "application/json")

	// Inizializza una variabile user e tenta di decodificare il corpo della richiesta in questa variabile.
	var nickname Nickname
	err := json.NewDecoder(r.Body).Decode(&nickname)

	// Controlla se c'è stato un errore durante la decodifica o se l'identificatore dell'utente non è valido.
	// In entrambi i casi, risponde con un codice di stato 400 Bad Request.
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !validNick(nickname.Nickname) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Crea l'utente nel database.
	// Se c'è un errore durante la creazione dell'utente nel database (ad es. l'utente esiste già),
	// risponde con un codice di stato 200 OK e restituisce l'utente.  Se c'è un errore durante la codifica della risposta,
	// risponde con un codice di stato 500 Internal Server Error.
	err = rt.db.CreateUser(nickname.Nickname)
	if err != nil {
		// l'utente esiste gia, viene restituito id e nickname
		id, err := rt.db.GetId(nickname.Nickname)
		if err != nil {
			// Caso di errore
			http.Error(w, "Can't get id.", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		u := User{
			UserId:   id,
			Nickname: nickname.Nickname,
		}
		err = json.NewEncoder(w).Encode(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("session: can't create response json")
		}
		return
	}

	// Se tutto va bene, risponde con un codice di stato 201 Created e restituisce l'utente
	// Se c'è un errore rispondo con error 500
	id, err := rt.db.GetId(nickname.Nickname)
	if err != nil {
		// Caso di errore
		http.Error(w, "Can't get id", http.StatusInternalServerError)
		return
	}
	u := User{
		UserId:   id,
		Nickname: nickname.Nickname,
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: can't create response json")
		return
	}
}
