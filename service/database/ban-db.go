package database

import (
// "database/sql"
// "errors"
// "fmt"
)

// Database fuction che permette a un utente(banner) di bannarne un'altro(banned)
func (db *appdbimpl) BanUser(banner int64, banned int64) error {

	_, err := db.c.Exec("INSERT INTO banned_users (banner,banned) VALUES (?, ?)", banner, banned)
	if err != nil {
		return err
	}

	return nil
}

// Database fuction che rimuovere un utente(banned)dalla lista dei banned di un'altro utente(banner)
func (db *appdbimpl) UnbanUser(banner int64, banned int64) error {

	_, err := db.c.Exec("DELETE FROM banned_users WHERE (banner = ? AND banned = ?)", banner, banned)
	if err != nil {
		return err
	}

	return nil
}

// Data base function per controllare se un utente è stato bannato.
// Restituisco 'true' se è banned, sennò 'false'
func (db *appdbimpl) BannedUserCheck(requestingUser int64, targetUser int64) (bool, error) {
	// Utilizza il metodo QueryRow per eseguire una query SQL SELECT COUNT(*) che
	// conta quante volte l'utente requestingUser appare nella tabella banned_users come utente bannato da targetUser
	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM banned_users WHERE banned = ? AND banner = ?",
		requestingUser, targetUser).Scan(&cnt)

	if err != nil {
		// Caso errore
		return true, err
	}

	// Se cnt e' pari a 1 esiste gia il ban
	if cnt == 1 {
		return true, nil
	}
	return false, nil
}
