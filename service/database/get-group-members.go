package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) GetGroupMembers(requesterID string, groupID int) ([]string, error) {

	// verifico che chi fa la richiesta sia membro del gruppo
	var isMember int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Ut_grup WHERE utente = ? AND gruppo = ?", requesterID, groupID).Scan(&isMember)
	if err != nil {
		return nil, fmt.Errorf("errore durante la verifica dell'appartenenza al gruppo: %w", err)
	}
	if isMember == 0 {
		return nil, errors.New("operazione negata: l'utente richiedente non è membro di questo gruppo")
	}

	// recupero degli username dei membri
	rows, err := db.c.Query(`
		SELECT u.username 
		FROM Utenti u
		INNER JOIN Ut_grup ug ON u.id = ug.utente
		WHERE ug.gruppo = ?
	`, groupID)
	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero dei membri del gruppo: %w", err)
	}
	defer rows.Close()

	// itero sui risultati per costruire la lista di utenti
	var members []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, fmt.Errorf("errore durante la scansione di un membro: %w", err)
		}
		members = append(members, username)
	}

	// controllo eventuali errori avvenuti durante l'iterazione di rows.Next()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// ritorno l'array vuoto o popolato
	if members == nil {
		members = make([]string, 0)
	}

	return members, nil
}
