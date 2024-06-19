package api

/*
definisce un insieme di endpoint per un'API web utilizzando un router HTTP.
Ogni endpoint è associato a un metodo HTTP specifico e a un "handler" che gestisce le richieste a quell'endpoint.
Handler serve a configurare un router HTTP che sa come gestire diverse richieste HTTP a vari endpoint dell'API web,
indirizzando ciascuna richiesta all'handler appropriato e assicurando che la logica dell'applicazione sia eseguita correttamente.
    Ricezione della Richiesta:
        Quando una richiesta HTTP arriva al server, il Handler determina quale funzione handler dovrebbe gestire la richiesta basata sull'URL e il metodo HTTP.
    Utilizzo di wrap:
        L'handler selezionato è "avvolto" dalla funzione wrap. La funzione wrap può eseguire del codice prima di chiamare
		l'handler principale, come configurare il RequestContext o eseguire il logging.
        wrap può anche eseguire del codice dopo che l'handler principale ha finito, come il logging aggiuntivo o la gestione degli errori.
    Configurazione del RequestContext:
        All'interno della funzione wrap, un RequestContext viene creato e configurato. Questo potrebbe includere la
		generazione di un ID univoco per la richiesta e la configurazione di un logger per includere quell'ID nelle voci di log.
    Chiamata dell'Handler Principale:
        La funzione wrap chiama l'handler principale, passando la richiesta HTTP originale e il RequestContext configurato.
    Elaborazione della Richiesta:
        L'handler principale elabora la richiesta, utilizzando le informazioni nel RequestContext come necessario.
		Ad esempio, potrebbe utilizzare il logger nel RequestContext per registrare messaggi di log che includono l'ID della richiesta.
    Risposta:
        L'handler principale genera una risposta HTTP e la restituisce al client attraverso il Handler.
*/

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
// Handler restituisce un oggetto che implementa l'interfaccia http.Handler, che può essere utilizzato per gestire le richieste HTTP.

func (rt *_router) Handler() http.Handler {

	// Login enpoint
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// Search endpoint
	rt.router.GET("/user", rt.wrap(rt.getUserProfiles))

	// Users Endpoint
	rt.router.GET("/user/:id", rt.wrap(rt.getUserProfile))
	rt.router.PUT("/user/:id/nickname", rt.wrap(rt.setMyUserName))
	rt.router.GET("/user/:id/nickname", rt.wrap(rt.getUserName))

	// Ban endpoint
	rt.router.PUT("/user/:id/banned_users", rt.wrap(rt.banUser))
	rt.router.DELETE("/user/:id/banned_users", rt.wrap(rt.unbanUser))

	// Followers endpoint
	rt.router.PUT("/user/:id/followers", rt.wrap(rt.followUser))
	rt.router.DELETE("/user/:id/followers", rt.wrap(rt.unfollowUser))

	// Stream endpoint
	rt.router.GET("/user/:id/stream", rt.wrap(rt.getMyStream))

	// Post Endpoint
	rt.router.POST("/user/:id/posts", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/user/:id/posts", rt.wrap(rt.deletePhoto))
	rt.router.GET("/user/:id/posts", rt.wrap(rt.getPhoto))

	// Comments endpoint
	rt.router.POST("/user/:id/posts/comments", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/user/:id/posts/comments", rt.wrap(rt.uncommentPhoto))

	// Likes endpoint
	rt.router.PUT("/user/:id/posts/likes", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/user/:id/posts/likes", rt.wrap(rt.unlikePhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
