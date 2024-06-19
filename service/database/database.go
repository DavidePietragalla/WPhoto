/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.
To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.
For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
Questo codice serve come middleware tra l'applicazione WASAPhoto e il suo database,fornendo
funzionalità per eseguire operazioni CRUD (Create, Read, Update, Delete) sul database in modo strutturato e organizzato.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Errors section
var ErrPostDoesntExist = errors.New("post doesn't exist")
var ErrUserBanned = errors.New("user is banned")

/*
var ErrUserAutoLike = errors.New("users can't like their own posts")
var ErrUserAutoFollow = errors.New("users can't follow themselfes")
*/

// Constants che indica foto per home
const PostsPerUserHome = 3

// AppDatabase è l'interfaccia per i DB
// L'interfaccia AppDatabase definisce un set di metodi che qualsiasi implementazione del database dovrebbe fornire.
type AppDatabase interface {

	// Crea un nuovo utente nel database ^
	CreateUser(string) error

	// Modifica il nickname di un utente ^
	ChangeNickname(int64, string) error

	// Restituisce gli utenti che matchano con il nick fornito ^
	SearchUser(string) ([]User, error)

	// Crea un nuovo post nel database ^
	CreatePost(int64, time.Time, []byte) (int64, error)

	// Recupera uno specifico post dal database ^
	GetPost(int64, int64) (Post, error)

	// Inserisce un like nel database ^
	LikePost(int64, int64) error

	// Rimuove un like dal database ^
	UnlikePost(int64, int64) error

	// Aggiunge un commento sotto un post ^
	CommentPost(int64, int64, string) (int64, error)

	// L'autore del commento decide di rimuoverlo ^
	UncommentPost(int64, int64) error

	// Aggiunge un follow al database ^
	FollowUser(int64, int64) error

	// Rimuove un follow dal database ^
	UnfollowUser(int64, int64) error

	// Aggiunge un ban nel database ^
	BanUser(int64, int64) error

	// Rimuove un ban dal database ^
	UnbanUser(int64, int64) error

	// Restituisce lo stream dell'utente ^
	GetStream(User) ([]Post, error)

	// Rimuove un post dal database ^
	RemovePost(int64, int64) error

	// Restituisce una lista di coloro che seguono l'utente ^
	GetFollowers(int64) ([]User, error)

	// Restituisce una lista di coloro che sono seguiti dall'utente ^
	GetFollowing(int64) ([]User, error)

	// Restituisce la lista di post di un utente a patto che il richiedente non sia bannato da esso ^
	GetPostsList(int64, int64) ([]Post, error)

	// Recupera l'id dato un nickname ^
	GetId(string) (int64, error)

	// Recupera il nickname dato un id ^
	GetNickname(int64) (string, error)

	// Recupera l'id dell'autore di un post ^
	GetPostAuthor(int64) (int64, error)

	// Recupera l'id dell'autore del commento ^
	GetCommentAuthor(int64) (int64, error)

	// Controlla se un ban gia esiste ^
	BannedUserCheck(int64, int64) (bool, error)

	// Controlla se un utente gia esiste ^
	CheckUser(string) (bool, error)

	// Ping checks whether the database is available or not (in that case, an error will be returned)
	Ping() error
}

// Questa struttura contiene un campo c che rappresenta la connessione al database.
type appdbimpl struct {
	c *sql.DB
}

// Questa funzione crea e restituisce una nuova istanza di AppDatabase. Se il database non esiste, viene creato.
func New(db *sql.DB) (AppDatabase, error) {
	// non è stata fornita una connessione valida al database.
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Esegue una query SQL per attivare il supporto alle chiavi esterne in SQLite.
	// Questo è importante per garantire l'integrità referenziale tra le tabelle.
	_, errPramga := db.Exec(`PRAGMA foreign_keys= ON`)
	if errPramga != nil {
		return nil, fmt.Errorf("error setting pragmas: %w", errPramga)
	}

	// Esegue una query SQL per verificare se esiste una tabella chiamata 'users' nel database.
	// Il risultato (il nome della tabella) viene memorizzato nella variabile tableName.
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		// Se la tabella 'users' non esiste, chiama la funzione createDatabase per creare tutte le tabelle necessarie nel database.
		err = createDatabase(db)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}
	// restituisce un'istanza dell'implementazione appdbimpl dell'interfaccia AppDatabase
	return &appdbimpl{
		c: db,
	}, nil
}

// Questa funzione verifica se il database è disponibile o meno
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// funzione crea tutte le tabelle necessarie per l'applicazione WASAPhoto nel database.
// Utilizza una serie di stringhe SQL per definire le tabelle e le relazioni tra di esse.
func createDatabase(db *sql.DB) error {
	// Ogni stringa rappresenta una query SQL per creare una tabella specifica nel database.
	// Le 6 stringhe comandi SQL per creare le tabelle users, posts, likes, comments, banned_users e followers se non esistono già.
	tables := [6]string{
		`CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			nickname VARCHAR(16) NOT NULL
			);`,
		`CREATE TABLE IF NOT EXISTS posts (
			post_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			date DATETIME NOT NULL,
			photo BLOB,
			FOREIGN KEY(user_id) REFERENCES users (user_id) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS  likes (
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			PRIMARY KEY (post_id, user_id),
			FOREIGN KEY(post_id) REFERENCES posts (post_id) ON DELETE CASCADE,
			FOREIGN KEY(user_id) REFERENCES users (user_id) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS comments (
			comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content VARCHAR(30) NOT NULL,
			FOREIGN KEY(post_id) REFERENCES posts (post_id) ON DELETE CASCADE,
			FOREIGN KEY(user_id) REFERENCES users (user_id) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS banned_users (
			banner INTEGER NOT NULL,
			banned INTEGER NOT NULL,
			PRIMARY KEY (banner,banned),
			FOREIGN KEY(banner) REFERENCES users (user_id) ON DELETE CASCADE,
			FOREIGN KEY(banned) REFERENCES users (user_id) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS followers(
			follower INTEGER NOT NULL,
			followed INTEGER NOT NULL,
			PRIMARY KEY (follower,followed),
			FOREIGN KEY(follower) REFERENCES users (user_id) ON DELETE CASCADE,
			FOREIGN KEY(followed) REFERENCES users (user_id) ON DELETE CASCADE
			);`,
	}

	// Inizia un ciclo for che itera su ogni stringa nell'array tables,per eseguire ogni query SQL nell'array per creare le tabelle.
	for i := 0; i < len(tables); i++ {
		// Assegno la query SQL corrente dall'array tables alla variabile sqlStmt.
		// utilizzo il metodo Exec,se la tabella specificata dalla query esiste già,la query non avrà effetto grazie al IF NOT EXISTS nelle stringhe SQL.
		sqlStmt := tables[i]
		_, err := db.Exec(sqlStmt)

		if err != nil {
			return err
		}
	}
	return nil
}
