package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) CreateGroup(identifier string, name string, members []string) (int, error) {
	// mappa per raccogliere gli ID unici dei partecipanti al gruppo dato che go non ha gli insiemi
	userIDsMap := make(map[string]bool)
	userIDsMap[identifier] = true // se l'identificatore non esiste viene inserito

	// traduzione degli username dei membri passati nel payload nei loro id utente
	for _, username := range members {
		var targetUserID string
		err := db.c.QueryRow("SELECT id FROM Utenti WHERE username = ?", username).Scan(&targetUserID)
		if err != nil {
			if err == sql.ErrNoRows {
				return 0, fmt.Errorf("utente membro non trovato: %s", username)
			}
			return 0, err
		}
		userIDsMap[targetUserID] = true
	}

	// avvio una transazione per garantire l'atomicitá dell'operazione
	tx, err := db.c.Begin()
	if err != nil {
		return 0, err
	}

	// creo la tupla base nella tabella Conversazioni per generare l'id condiviso
	res, err := tx.Exec("INSERT INTO Conversazioni DEFAULT VALUES")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	groupID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// creo il gruppo
	_, err = tx.Exec("INSERT INTO Gruppi (id, nome) VALUES (?, ?)", groupID, name)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// associo i membri al gruppo
	for userID := range userIDsMap {
		_, err = tx.Exec("INSERT INTO Ut_grup (utente, gruppo) VALUES (?, ?)", userID, groupID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// se non ci sono errori, si procede
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(groupID), nil
}
