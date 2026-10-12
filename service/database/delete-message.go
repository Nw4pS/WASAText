package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// requesterID è l'id dell'utente che effettua la richiesta (estratto dal token)
// messageID è l'id del messaggio da eliminare
func (db *appdbimpl) DeleteMessage(requesterID string, messageID int) error {
	// recupero le informazioni del messaggio per verifica permessi e stato
	var senderID string
	var isDeleted int

	err := db.c.QueryRow("SELECT utente, is_eliminato FROM Messaggi WHERE id = ?", messageID).Scan(&senderID, &isDeleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("messaggio non trovato")
		}
		return fmt.Errorf("errore durante la ricerca del messaggio: %w", err)
	}

	// controllo che il messaggio non sia già eliminato
	if isDeleted == 1 {
		return nil
	}

	// controllo se l'utente può eliminare il messaggio
	if senderID != requesterID {
		return errors.New("operazione negata: non puoi eliminare un messaggio inviato da un altro utente")
	}

	// soft delete, aggiorno il flag is_eliminato
	res, err := db.c.Exec("UPDATE Messaggi SET is_eliminato = 1 WHERE id = ?", messageID)
	if err != nil {
		return fmt.Errorf("errore durante l'eliminazione logica del messaggio: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("errore nel controllo delle righe eliminate: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("impossibile aggiornare lo stato del messaggio")
	}

	return nil
}
