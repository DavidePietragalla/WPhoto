package database

import (
// "database/sql"
// "errors"
// "fmt"
)

// Database function che  permette a un utente di mettere "mi piace" a una foto.
func (db *appdbimpl) LikePost(postId int64, u int64) error {
	// Aggiunge un record nella tabella likes, indicando che l'utente u ha messo "mi piace" alla foto.
	_, err := db.c.Exec("INSERT INTO likes (post_id, user_id) VALUES (?, ?)", postId, u)
	if err != nil {
		return err
	}

	return nil
}

// Database function che permette a un utente (u) di rimuovere il "mi piace" da una foto (p).
func (db *appdbimpl) UnlikePost(postId int64, u int64) error {
	// Rimuove un record dalla tabella likes, indicando che l'utente u ha rimosso il "mi piace" dalla foto.
	_, err := db.c.Exec("DELETE FROM likes WHERE(post_id = ? AND user_id = ?)", postId, u)
	if err != nil {
		return err
	}

	return nil
}

// Database function che recupera la lista degli utenti che hanno messo "mi piace" a una determinata foto.
func (db *appdbimpl) GetLikesList(postId int64) ([]User, error) {
	// La query SQL seleziona tutti gli utenti che hanno messo "mi piace" alla foto specificata
	rows, err := db.c.Query("SELECT user_id FROM likes WHERE post_id = ?", postId)
	if err != nil {
		return nil, err
	}
	// Attende che la funzione finisca per chiudere le rows
	defer func() { _ = rows.Close() }()

	// Legge gli user id nelle righe
	var likes []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserId)
		if err != nil {
			return nil, err
		}

		// Get the nickname of the user that liked the post
		nickname, err := db.GetNickname(user.UserId)
		if err != nil {
			return nil, err
		}
		user.Nickname = nickname

		likes = append(likes, user)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return likes, nil
}
