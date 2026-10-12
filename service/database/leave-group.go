package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) LeaveGroup(userID string, groupID int) error {
	res, err := db.c.Exec("DELETE FROM Ut_grup WHERE utente = ? AND gruppo = ?", userID, groupID)
	if err != nil {
		return fmt.Errorf("errore durante la query di uscita dal gruppo: %w", err)
	}

	// verifico quante righe sono state effettivamente eliminate
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("errore nel controllo delle righe eliminate: %w", err)
	}

	// se nessuna riga è stata eliminata, l'utente non faceva parte del gruppo
	if rowsAffected == 0 {
		return errors.New("l'utente non fa parte di questo gruppo")
	}

	return nil
}
