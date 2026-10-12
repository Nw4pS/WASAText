package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) AddReaction(requesterID string, messageID int, emoji string) error {
	// id a della conversazione a cui appartiene il messaggio
	var convID int
	err := db.c.QueryRow("SELECT conversazione FROM Messaggi WHERE id = ?", messageID).Scan(&convID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("messaggio non trovato")
		}
		return fmt.Errorf("errore durante la ricerca del messaggio: %w", err)
	}

	// controllo che l'utente abbia accesso alla conversazione
	var isMember int
	// controllo se l'id della conversazione è conversazione diretta
	err = db.c.QueryRow("SELECT COUNT(*) FROM Ut_dir WHERE utente = ? AND diretto = ?", requesterID, convID).Scan(&isMember)
	if err != nil {
		return fmt.Errorf("errore nel controllo appartenenza chat diretta: %w", err)
	}

	if isMember == 0 { // controllo se è id del gruppo
		err = db.c.QueryRow("SELECT COUNT(*) FROM Ut_grup WHERE utente = ? AND gruppo = ?", requesterID, convID).Scan(&isMember)
		if err != nil {
			return fmt.Errorf("errore nel controllo appartenenza gruppo: %w", err)
		}
	}

	if isMember == 0 {
		return errors.New("operazione negata: l'utente non ha accesso a questo messaggio")
	}

	// controllo se la reazione è già presente
	var existingReactionID int
	err = db.c.QueryRow("SELECT id FROM Reaction WHERE utente = ? AND messaggio = ?", requesterID, messageID).Scan(&existingReactionID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// se non c'è reazione dall'utente, la inserisce
			_, err = db.c.Exec("INSERT INTO Reaction (utente, Messaggio, contenuto) VALUES (?, ?, ?)", requesterID, messageID, emoji)
			if err != nil {
				return fmt.Errorf("errore durante l'inserimento della reazione: %w", err)
			}
			return nil
		}
		return fmt.Errorf("errore durante la verifica delle reazioni: %w", err)
	}

	// se la reazione esiste, la si aggiorna
	_, err = db.c.Exec("UPDATE Reaction SET contenuto = ? WHERE id = ?", emoji, existingReactionID)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamento della reazione: %w", err)
	}

	return nil
}
