package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

// creo nuovi errori per poter poterli passare all'handler e distinguerli meglio
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUsernameTaken = errors.New("username already in use")
)

func (db *appdbimpl) SetMyUsername(requesterID string, newUsername string) error {
	// verifico se l'utente esiste e recupero lo username corrente
	var currentUsername string
	err := db.c.QueryRow("SELECT username FROM Utenti WHERE id = ?", requesterID).Scan(&currentUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("errore durante la ricerca dell'utente: %w", err)
	}

	// se l'utente invia lo stesso username che ha già
	if currentUsername == newUsername {
		return nil
	}

	// eseguo l'update
	_, err = db.c.Exec("UPDATE Utenti SET username = ? WHERE id = ?", newUsername, requesterID)
	if err != nil {
		// verifico se l'errore è una violazione del vincolo UNIQUE su username
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return ErrUsernameTaken
		}
		return fmt.Errorf("errore durante l'aggiornamento dell'username: %w", err)
	}

	return nil
}
