package database

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrGroupNotFound  = errors.New("gruppo non trovato")
	ErrNotGroupMember = errors.New("operazione negata: l'utente non fa parte del gruppo")
	ErrInvalidName    = errors.New("nome del gruppo non valido")
)

func (db *appdbimpl) SetGroupName(requesterID string, groupID int, newName string) error {
	trimmedName := strings.TrimSpace(newName)
	if trimmedName == "" {
		return ErrInvalidName
	}

	// controllo che il richiedente sia essere membro del gruppo
	var isMember int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Ut_grup WHERE utente = ? AND gruppo = ?", requesterID, groupID).Scan(&isMember)
	if err != nil {
		return fmt.Errorf("errore verifica appartenenza gruppo: %w", err)
	}
	if isMember == 0 {
		return ErrNotGroupMember
	}

	// aggiorno il nome nella tabella Gruppi
	res, err := db.c.Exec("UPDATE Gruppi SET nome = ? WHERE id = ?", trimmedName, groupID)
	if err != nil {
		return fmt.Errorf("errore aggiornamento nome gruppo: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("errore verifica righe aggiornate: %w", err)
	}
	if rowsAffected == 0 {
		return ErrGroupNotFound
	}

	return nil
}
