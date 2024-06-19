package database

import ()

type User struct {
	UserId   int64  `json:"user_id"`  // Id univoco dell'utente
	Nickname string `json:"nickname"` // Nickname dell'utente
}

// Database function che aggiunge un nuovo utente al database durante la registrazione.
func (db *appdbimpl) CreateUser(u string) error {
	// Controllo che il nickname non sia già stato preso
	check, _ := db.CheckUser(u)
	if check {
		return nil
	}

	// Esegue una query SQL per inserire un nuovo utente nel database con un ID utente e un nickname
	_, err := db.c.Exec("INSERT INTO users (nickname) VALUES (?)", u)

	if err != nil {
		return err
	}

	return nil
}

// Database function che restituice l'id di un utente
func (db *appdbimpl) GetId(nickname string) (int64, error) {

	var id int64

	// Utilizza una query SQL SELECT per cercare l'id dell'utente
	err := db.c.QueryRow(`SELECT user_id FROM users WHERE nickname = ?`, nickname).Scan(&id)
	if err != nil {
		// Caso di errore
		return -1, err
	}
	return id, nil
}

// Database function che restituice il nickname di un utente
func (db *appdbimpl) GetNickname(userId int64) (string, error) {

	var nickname string

	// Utilizza una query SQL SELECT per cercare il nickname dell'utente
	err := db.c.QueryRow(`SELECT nickname FROM users WHERE user_id = ?`, userId).Scan(&nickname)
	if err != nil {
		// Caso di errore
		return "", err
	}
	return nickname, nil
}

// Modifica il nickname di un utente
func (db *appdbimpl) ChangeNickname(u int64, newNick string) error {
	// Controllo che il nuovo nickname non sia già occupato
	check, _ := db.CheckUser(newNick)
	if check {
		return nil
	}

	// Sostituisco il nickname
	_, err := db.c.Exec(`UPDATE users SET nickname = ? WHERE user_id = ?`, newNick, u)

	if err != nil {
		return err
	}

	return nil
}

// Database function controlla se un utente esiste nel database.
func (db *appdbimpl) CheckUser(n string) (bool, error) {
	//  Esegue una query SQL per contare il numero di righe nella tabella degli utenti che corrispondono all'ID utente specificato.
	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE nickname = ?",
		n).Scan(&cnt)

	if err != nil {
		// Nel caso di un errore
		return true, err
	}

	// Se il contatore è 1 l'utente esiste e restituisco True
	if cnt == 1 {
		return true, nil
	}
	return false, nil
}

// Database function che permette di cercare utenti in base a un nickname fornito.
// Ogni partial macth viene incluso nei risultati,restituendo una lista di macthing users
func (db *appdbimpl) SearchUser(userToSearch string) ([]User, error) {
	// La funzione esegue una query SQL SELECT per cercare tutti gli utenti il cui nickname corrisponde parzialmente
	// al parametro fornito. L'uso del simbolo % dopo userToSearch nella query indica una ricerca di tipo "LIKE",
	// che restituirà tutti gli utenti che hanno un nickname che inizia con il valore di userToSearch.
	rows, err := db.c.Query("SELECT * FROM users WHERE (nickname LIKE ?)",
		userToSearch+"%")
	if err != nil {
		return nil, err
	}
	// Attende che la funzione finisca per chiudere le rows
	defer func() { _ = rows.Close() }()

	// Inserisce i risultati dentro una lista
	var res []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserId, &user.Nickname)
		if err != nil {
			return nil, err
		}
		res = append(res, user)
	}
	// Restituisce un elenco di utenti (User) che corrispondono al parametro di ricerca fornito.
	// Se si verifica un errore durante l'esecuzione della query o la lettura dei risultati, viene restituito un errore.
	if rows.Err() != nil {
		return nil, err
	}

	return res, nil
}

// Database function che recupera lo "stream" di un utente, che consiste nelle foto delle persone seguite dall'utente.
func (db *appdbimpl) GetStream(user User) ([]Post, error) {
	// Esegue una query SQL SELECT per selezionare tutte le foto degli utenti seguiti
	// dall'utente specificato e le ordina in base alla data di caricamento in ordine decrescente.
	rows, err := db.c.Query(`SELECT * FROM posts WHERE user_id IN (SELECT followed FROM followers WHERE follower = ?) ORDER BY date DESC`,
		user.UserId)
	if err != nil {
		return nil, err
	}
	// Attende che la funzione finisca per chiudere le rows
	defer func() { _ = rows.Close() }()

	// Salva tutti i post in un array
	var res []Post
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.PostId, &post.UserId, &post.Date, &post.Url)
		if err != nil {
			return nil, err
		}
		res = append(res, post)
	}

	if rows.Err() != nil {
		return nil, err
	}
	// Restituisce una slice di Post che rappresenta lo stream dell'utente.
	return res, nil
}
