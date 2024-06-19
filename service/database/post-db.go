package database

import (
	// "database/sql"
	// "errors"
	"time"
)

type Post struct {
	PostId   int64     `json:"post_id"`  // Unique id of the post
	UserId   int64     `json:"owner"`    // Unique id of the owner
	Date     time.Time `json:"date"`     // Date in which the post was uploaded
	Url      []byte    `json:"urlImage"` // File of the picture
	Comments []Comment `json:"comments"` // Array of comments of the post
	Likes    []User    `json:"likes"`    // Array of useres that liked the post
}

// Database function che crea una nuova foto nel database e restituisce l'ID univoco della foto.
func (db *appdbimpl) CreatePost(userId int64, date time.Time, url []byte) (int64, error) {
	// Utilizza una query SQL INSERT per inserire la foto nel database.
	res, err := db.c.Exec("INSERT INTO posts (user_id, date, photo) VALUES (?,?,?)",
		userId, date, url)

	if err != nil {
		// Caso di errore
		return -1, err
	}

	postId, err := res.LastInsertId()
	if err != nil {
		// Errore prendendo l'id dell'ultima foto inserita (postId)
		return -1, err
	}

	return postId, nil
}

// Database function che recupera una foto specifica (targetPost),
// ma solo se l'utente che fa la richiesta (requestinUser) non è stato bannato dal proprietario della foto.
func (db *appdbimpl) GetPost(requestinUser int64, postId int64) (Post, error) {
	// Utilizza una query SQL SELECT per recuperare la foto nel database.
	var id int64
	var usid int64
	var date time.Time
	var url []byte
	err := db.c.QueryRow("SELECT * FROM posts WHERE (post_id = ?) AND user_id NOT IN (SELECT banner FROM banned_users WHERE banned = ?)",
		postId, requestinUser).Scan(&id, &usid, &date, &url)

	if err != nil {
		return Post{}, err
	}
	comments, _ := db.GetCompleteCommentsList(id, usid)
	likes, _ := db.GetLikesList(id)
	return Post{
		PostId:   id,
		UserId:   usid,
		Date:     date,
		Url:      url,
		Comments: comments,
		Likes:    likes,
	}, nil

}

// Rimuove una foto specifica dal database, ma solo se l'utente specificato (owner) è il proprietario della foto.
func (db *appdbimpl) RemovePost(requestinUser int64, postId int64) error {

	_, err := db.c.Exec("DELETE FROM posts WHERE user_id = ? AND post_id = ? ",
		requestinUser, postId)
	if err != nil {
		// In caso di errore
		return err
	}

	return nil
}

func (db *appdbimpl) GetPostAuthor(postId int64) (int64, error) {
	// Utilizza una query SQL SELECT per recuperare la foto nel database.
	var userId int64
	err := db.c.QueryRow("SELECT  user_id FROM posts WHERE (post_id = ?)", postId).Scan(&userId)

	if err != nil {
		return -1, err
	}

	return userId, nil

}

// Database function che recupera l'elenco delle foto di un utente (targetUser),
// ma solo se l'utente che fa la richiesta (requestingUser) non è stato bannato da targetUser.
func (db *appdbimpl) GetPostsList(requestingUser int64, targetUser int64) ([]Post, error) {
	// Controlla se requestingUser è bannato da targetUser
	check, err := db.BannedUserCheck(requestingUser, targetUser)
	if err != nil {
		return nil, err
	}
	if check {
		return nil, err
	}

	// Esegue una query SQL per selezionare tutte le foto di targetUser e le ordina in base alla data in ordine decrescente.
	rows, err := db.c.Query("SELECT * FROM posts WHERE user_id = ? ORDER BY date DESC", targetUser)
	if err != nil {
		return nil, err
	}
	// Attende che la funzione finisca per chiudere le rows
	defer func() { _ = rows.Close() }()

	// Recupera il contenuto di ogni post e lo inserisce in una lista
	var posts []Post
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.PostId, &post.UserId, &post.Date, &post.Url)
		if err != nil {
			return nil, err
		}

		comments, err := db.GetCompleteCommentsList(post.PostId, post.UserId)
		if err != nil {
			return nil, err
		}
		post.Comments = comments

		likes, err := db.GetLikesList(post.PostId)
		if err != nil {
			return nil, err
		}
		post.Likes = likes

		posts = append(posts, post)
	}

	if rows.Err() != nil {
		return nil, err
	}
	// Restituisce un elenco di foto con i relativi commenti e "mi piace".
	return posts, nil
}
