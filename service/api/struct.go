package api

import (
	"WasaPhotoDavidePietragalla/service/database"
)

// Qui vengono definiti vari tipi di dati e metodi per rappresentare e gestire utenti,foto,commenti
// e altri dati correlati nell'API. Include anche metodi per convertire questi tipi di dati in tipi di dati utilizzati dal database.
// Dichiarazione dei Error messages
const INTERNAL_ERROR_MSG = "internal server error"
const PNG_ERROR_MSG = "file is not a png format"
const JPG_ERROR_MSG = "file is not a jpg format"
const IMG_FORMAT_ERROR_MSG = "images must be jpeg or png"
const INVALID_JSON_ERROR_MSG = "invalid json format"
const INVALID_IDENTIFIER_ERROR_MSG = "identifier must be a string between 3 and 16 characters"

// JSON Error Structure : Definizione di una struttura per rappresentare un messaggio di errore in formato JSON.
type JSONErrorMsg struct {
	Message string `json:"message"` // Error messages
}

// User structure for the APIs
type User struct {
	UserId   int64  `json:"user_id"`  // User's unique id
	Nickname string `json:"nickname"` // Nickname of a user
}

// Nickname structure for the APIs
type Nickname struct {
	Nickname string `json:"nickname"` // Nickname of a user
}

// URL structure for the APIs
type Url struct {
	Data []byte `json:"urlImage"` // Url of the image
}

// Comment structure for the APIs
type Comment struct {
	Comment string `json:"comment"` // Comment content
}

// CommentId structure for the APIs
type CommentId struct {
	CommentId int64 `json:"comment_id"` // Identifier of a comment
}

// CompleteProfile structure for the APIs
type Profile struct {
	UserId    int64           `json:"user_id"`
	Nickname  string          `json:"nickname"`
	Followers []database.User `json:"followers"`
	Following []database.User `json:"following"`
	Posts     []database.Post `json:"posts"`
}
