package database

import (
// "database/sql"
// "errors"
// "fmt"
)

// Database function che permette a un utente (follower) di seguire un altro utente (followed).
func (db *appdbimpl) FollowUser(follower int64, followed int64) error {
	// Aggiunge un record nella tabella followers, indicando che l'utente follower segue l'utente followed.
	_, err := db.c.Exec("INSERT INTO followers (follower,followed) VALUES (?, ?)",
		follower, followed)
	if err != nil {
		return err
	}

	return nil
}

// Database function che permette a un utente (follower) di smettere di seguire un altro utente (followed).
func (db *appdbimpl) UnfollowUser(follower int64, followed int64) error {
	// Rimuove un record dalla tabella followers,
	// indicando che l'utente follower non segue più l'utente followed.
	_, err := db.c.Exec("DELETE FROM followers WHERE(follower = ? AND followed = ?)",
		follower, followed)
	if err != nil {
		return err
	}

	return nil
}

// Database function che recupera la lista degli utenti che seguono l'utente specificato
func (db *appdbimpl) GetFollowers(userId int64) ([]User, error) {
	// Utilizza una query SQL per selezionare tutti gli utenti che seguono l'utente specificato.
	rows, err := db.c.Query("SELECT follower FROM followers WHERE followed = ?", userId)
	if err != nil {
		return nil, err
	}
	// Chiude le rows una volta finita la funzione
	defer func() { _ = rows.Close() }()

	// Inserisco in una lista gli utenti trovati
	var followers []User
	for rows.Next() {
		var follower User
		err = rows.Scan(&follower.UserId)
		if err != nil {
			return nil, err
		}

		nickname, err := db.GetNickname(follower.UserId)
		if err != nil {
			return nil, err
		}
		follower.Nickname = nickname

		followers = append(followers, follower)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return followers, nil
}

// Database function che recupera la lista degli utenti seguiti dall'utente specificato.
func (db *appdbimpl) GetFollowing(userId int64) ([]User, error) {
	// Utilizza una query SQL per selezionare tutti gli utenti seguiti dall'utente specificato.
	rows, err := db.c.Query("SELECT followed FROM followers WHERE follower = ?", userId)
	if err != nil {
		return nil, err
	}
	// Chiude le rows una volta finita la funzione
	defer func() { _ = rows.Close() }()

	// Inserisco in una lista gli utenti trovati
	var following []User
	for rows.Next() {
		var followed User
		err = rows.Scan(&followed.UserId)
		if err != nil {
			return nil, err
		}

		nickname, err := db.GetNickname(followed.UserId)
		if err != nil {
			return nil, err
		}
		followed.Nickname = nickname

		following = append(following, followed)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return following, nil
}
