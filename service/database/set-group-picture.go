package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) SetGroupPhoto(requesterID string, groupID int, photo string) error {
	// verifico che il richiedente sia membro del gruppo
	var isMember int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Ut_grup WHERE utente = ? AND gruppo = ?", requesterID, groupID).Scan(&isMember)
	if err != nil {
		return fmt.Errorf("errore durante la verifica dell'appartenenza al gruppo: %w", err)
	}
	if isMember == 0 {
		return errors.New("operazione negata: l'utente non è membro di questo gruppo")
	}

	// aggiorno la foto
	res, err := db.c.Exec("UPDATE Gruppi SET groupPic = ? WHERE id = ?", photo, groupID)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamento della foto del gruppo: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("errore nel controllo delle righe modificate: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("gruppo non trovato")
	}

	return nil
}
