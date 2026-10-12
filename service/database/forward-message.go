package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (db *appdbimpl) ForwardMessage(requesterID string, destConvID int, origMessageID int) (int, error) {
	// controllo se l'utente ha accesso alla conversazione di destinazione
	var isDestMember int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Ut_dir WHERE utente = ? AND diretto = ?", requesterID, destConvID).Scan(&isDestMember)
	if err != nil {
		return 0, fmt.Errorf("errore controllo appartenenza chat diretta destinazione: %w", err)
	}

	if isDestMember == 0 {
		err = db.c.QueryRow("SELECT COUNT(*) FROM Ut_grup WHERE utente = ? AND gruppo = ?", requesterID, destConvID).Scan(&isDestMember)
		if err != nil {
			return 0, fmt.Errorf("errore controllo appartenenza gruppo destinazione: %w", err)
		}
	}

	if isDestMember == 0 {
		return 0, errors.New("operazione negata: l'utente non fa parte della conversazione di destinazione")
	}

	// recupero il messaggio originale
	var origContent string
	var origIsImage int
	var origIsDeleted int
	var origConvID int

	err = db.c.QueryRow(`
		SELECT contenuto, is_immagine, is_eliminato, conversazione 
		FROM Messaggi 
		WHERE id = ?
	`, origMessageID).Scan(&origContent, &origIsImage, &origIsDeleted, &origConvID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("messaggio originale non trovato")
		}
		return 0, fmt.Errorf("errore lettura messaggio originale: %w", err)
	}

	if origIsDeleted == 1 {
		return 0, errors.New("operazione negata: non puoi inoltrare un messaggio eliminato")
	}

	// controllo se l'utente ha accesso alla conversazione di origine
	var isOrigMember int
	err = db.c.QueryRow("SELECT COUNT(*) FROM Ut_dir WHERE utente = ? AND diretto = ?", requesterID, origConvID).Scan(&isOrigMember)
	if err != nil {
		return 0, fmt.Errorf("errore controllo appartenenza chat diretta origine: %w", err)
	}

	if isOrigMember == 0 {
		err = db.c.QueryRow("SELECT COUNT(*) FROM Ut_grup WHERE utente = ? AND gruppo = ?", requesterID, origConvID).Scan(&isOrigMember)
		if err != nil {
			return 0, fmt.Errorf("errore controllo appartenenza gruppo origine: %w", err)
		}
	}

	if isOrigMember == 0 {
		return 0, errors.New("operazione negata: l'utente non ha accesso al messaggio originale")
	}

	// inserisco nel db il messaggio inoltrato
	currentTime := time.Now().UTC().Format(time.RFC3339)
	res, err := db.c.Exec(
		"INSERT INTO Messaggi (utente, conversazione, istante, contenuto, is_immagine, is_inoltrato, is_eliminato, stato) VALUES (?, ?, ?, ?, ?, 1, 0, 'sent')",
		requesterID, destConvID, currentTime, origContent, origIsImage,
	)
	if err != nil {
		return 0, fmt.Errorf("errore durante l'inserimento del messaggio inoltrato: %w", err)
	}

	newMsgID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("errore nel recupero dell'ID del nuovo messaggio: %w", err)
	}

	return int(newMsgID), nil
}
