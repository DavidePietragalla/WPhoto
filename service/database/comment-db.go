package database

import (
// "database/sql"
// "errors"
// "fmt"
)

type Comment struct {
	CommentId int64  `json:"comment_id"` // Identifier of a comment
	PostId    int64  `json:"post_id"`    // Post unique id
	UserId    int64  `json:"user_id"`    // User's unique id
	Nickname  string `json:"nickname"`   // Nickname of a user
	Content   string `json:"comment"`    // Comment content
}

// Database function per aggiungere un commento di un user ad una foto
func (db *appdbimpl) CommentPost(postId int64, u int64, content string) (int64, error) {
	// Utilizza una query SQL INSERT per inserire il commento nel database.
	res, err := db.c.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)",
		postId, u, content)
	if err != nil {
		// In caso di errore
		return -1, err
	}

	commentId, err := res.LastInsertId()
	if err != nil {
		// Errore nel recuperare l'id del commento
		return -1, err
	}
	// Restituisce l'ID del commento appena inserito.
	return commentId, nil
}

// Database function che rimuove il commento di un utente dalla foto
func (db *appdbimpl) UncommentPost(u int64, commentId int64) error {
	// Utilizza una query SQL DELETE per rimuovere il commento specificato dal database.
	_, err := db.c.Exec("DELETE FROM comments WHERE (comment_id = ? AND user_id = ?)",
		commentId, u)
	if err != nil {
		return err
	}

	return nil
}

// Database function che rimuove il commento di un utente dalla foto
func (db *appdbimpl) GetCommentAuthor(commentId int64) (int64, error) {
	// Utilizza una query SQL DELETE per trovare l'autore del commento specificato
	var authorId int64
	err := db.c.QueryRow("SELECT user_id FROM comments WHERE (comment_id = ?)", commentId).Scan(&authorId)
	if err != nil {
		return -1, err
	}

	return authorId, nil
}

// Questa funzione recupera la lista completa dei commenti di una foto, escludendo i commenti dei bannati
func (db *appdbimpl) GetCompleteCommentsList(postId int64, ownerId int64) ([]Comment, error) {
	// Utilizza una query SQL per selezionare tutti i commenti della foto specificata,
	// escludendo quelli degli utenti che sono bannati dal proprietario.
	rows, err := db.c.Query("SELECT * FROM comments WHERE post_id = ? AND user_id NOT IN (SELECT banned FROM banned_users WHERE banner = ?)",
		postId, ownerId)
	if err != nil {
		return nil, err
	}

	// Attende che la funzione finisca per chiudere le rows
	defer func() { _ = rows.Close() }()

	// Read all the comments in the resulset (comments of the post with authors that didn't ban the requesting user).
	var comments []Comment
	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.CommentId, &comment.PostId, &comment.UserId, &comment.Content)
		if err != nil {
			return nil, err
		}

		// Prende il nick dell'utente che ha commentato
		nickname, err := db.GetNickname(comment.UserId)
		if err != nil {
			return nil, err
		}
		comment.Nickname = nickname

		comments = append(comments, comment)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return comments, nil
}
