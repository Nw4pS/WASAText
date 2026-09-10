package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

// DoLogin cerca un utente tramite username, restituisce il suo ID se esiste
// altrimenti crea un nuovo utente generando un UUIDv4 e restituisce quello
func (db *appdbimpl) DoLogin(username string) (string, error) {
	var userID string

	// cerca l'utente nel database tramite il suo username univoco
	err := db.c.QueryRow("SELECT id FROM utenti WHERE username = ?", username).Scan(&userID)

	if err == nil {
		// l'utente esiste, quindi restituisce il suo identificatore
		return userID, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		// ritorna un errore se si verifica un problema diverso da "utente non trovato".
		return "", fmt.Errorf("errore durante la ricerca dell'utente: %w", err)
	}

	//l'utente non esiste, quindi genera un nuovo UUIDv4 per l'identificatore
	u, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("errore nella generazione dell'UUID: %w", err)
	}
	userID = u.String()

	// inserisce il nuovo utente con foto profilo inizializzata vuota.
	_, err = db.c.Exec("INSERT INTO utenti (id, username, proPic) VALUES (?, ?, '')", userID, username)
	if err != nil {
		return "", fmt.Errorf("errore durante la creazione dell'utente: %w", err)
	}

	return userID, nil
}
