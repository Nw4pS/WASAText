package database

import (
	"errors"
	"fmt"
)

var ErrReactionNotFound = errors.New("reazione non trovata")

func (db *appdbimpl) DeleteReaction(requesterID string, messageID int) error {
	res, err := db.c.Exec("DELETE FROM Reaction WHERE utente = ? AND Messaggio = ?", requesterID, messageID)
	if err != nil {
		return fmt.Errorf("errore cancellazione reazione: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrReactionNotFound
	}

	return nil
}
