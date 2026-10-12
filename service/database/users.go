package database

import (
	"database/sql"
)

type User struct {
	Username string `json:"username"`
	ProPic   string `json:"proPic"`
}

func (db *appdbimpl) SearchUsers(username string) ([]User, error) {
	var rows *sql.Rows
	var err error

	if username == "" {
		rows, err = db.c.Query("SELECT username, COALESCE(proPic,'') FROM utenti ORDER BY username ASC")
	} else {
		rows, err = db.c.Query("SELECT username, COALESCE(proPic,'') FROM utenti WHERE username = ?", username)
	}

	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	users := make([]User, 0)

	for rows.Next() {
		var u User

		if err := rows.Scan(&u.Username, &u.ProPic); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
