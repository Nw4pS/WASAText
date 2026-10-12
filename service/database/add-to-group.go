package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) AddToGroup(requesterID string, groupID int, targetUsernames []string) error {
	// verifico che chi fa la richiesta sia già membro del gruppo
	var isMember int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Ut_grup WHERE utente = ? AND gruppo = ?", requesterID, groupID).Scan(&isMember)
	if err != nil {
		return fmt.Errorf("errore durante la verifica dell'appartenenza al gruppo: %w", err)
	}
	if isMember == 0 {
		return errors.New("operazione negata: l'utente richiedente non è membro di questo gruppo")
	}

	// rimuovo eventuali duplicati dall'input usando una mappa (in go non esistono i set)
	uniqueUsernames := make(map[string]bool)
	for _, username := range targetUsernames {
		uniqueUsernames[username] = true
	}

	// inizio della transazione con il db
	tx, err := db.c.Begin()
	if err != nil {
		return fmt.Errorf("errore nell'avvio della transazione: %w", err)
	}

	// ciclo sugli username
	for username := range uniqueUsernames {
		// recupero l'id associato a quell'username
		var targetUserID string
		err = tx.QueryRow("SELECT id FROM Utenti WHERE username = ?", username).Scan(&targetUserID)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("utente membro non trovato: %s", username)
			}
			return fmt.Errorf("errore nella ricerca dell'utente %s: %w", username, err)
		}

		// verifico se l'utente fa già parte del gruppo.
		var alreadyInGroup int
		err = tx.QueryRow("SELECT COUNT(*) FROM Ut_grup WHERE utente = ? AND gruppo = ?", targetUserID, groupID).Scan(&alreadyInGroup)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("errore nel controllo appartenenza per %s: %w", username, err)
		}
		if alreadyInGroup > 0 {
			tx.Rollback()
			return fmt.Errorf("l'utente %s fa già parte di questo gruppo", username)
		}

		// si inserisce l'utente nel gruppo
		_, err = tx.Exec("INSERT INTO Ut_grup (utente, gruppo) VALUES (?, ?)", targetUserID, groupID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("errore durante l'aggiunta dell'utente %s al gruppo: %w", username, err)
		}
	}

	// se non ci sono errori si conferma la transazione
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("errore nel commit della transazione: %w", err)
	}

	return nil
}
