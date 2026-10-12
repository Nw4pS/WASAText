package database

import "database/sql"

func (db *appdbimpl) CreateConversation(identifier string, username string) (int, error) {

	var targetUserID string //id dell'utente a con cui si vuole creare la conversazione diretta

	err := db.c.QueryRow("SELECT id FROM Utenti WHERE Utenti.username=?", username).Scan(&targetUserID)
	if err != nil {
		return 0, err
	}

	var idChatDiretto int
	// controllo che la chat non esista già
	err = db.c.QueryRow(`SELECT ud1.diretto FROM Ut_dir ud1 
							JOIN Ut_dir ud2 ON ud1.diretto=ud2.diretto
							WHERE ud1.utente=? AND ud2.utente=?`, identifier, targetUserID).Scan(&idChatDiretto)

	if err == nil {
		// la chat esiste, ritorno l'id
		return idChatDiretto, nil
	}

	// in caso non esista:
	if err == sql.ErrNoRows {
		// inizio transazione con il db in modo tale da poter fare il rollback in caso di errori
		tx, err := db.c.Begin()
		if err != nil {
			return 0, err
		}

		// creo la tupla in Conversazioni per generare l'id
		res, err := tx.Exec("INSERT INTO Conversazioni DEFAULT VALUES")
		if err != nil {
			tx.Rollback() // rollback in caso di errore
			return 0, err
		}

		// recupero l'id appena creato
		convID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback() // rollback in caso di errore
			return 0, err
		}

		// inserisco il nuovo id nella tabella Diretti
		_, err = tx.Exec("INSERT INTO Diretti (id) VALUES (?)", convID)
		if err != nil {
			tx.Rollback() // rollback in caso di errore
			return 0, err
		}

		// collego il primo utente nella tabella Ut_dir
		_, err = tx.Exec("INSERT INTO Ut_dir (utente, diretto) VALUES (?, ?)", identifier, convID)
		if err != nil {
			tx.Rollback() // rollback in caso di errore
			return 0, err
		}

		// collego il secondo utente
		_, err = tx.Exec("INSERT INTO Ut_dir (utente, diretto) VALUES (?, ?)", targetUserID, convID)
		if err != nil {
			tx.Rollback() // Rollback in caso di errore
			return 0, err
		}

		// fine transazione
		if err = tx.Commit(); err != nil {
			return 0, err
		}

		// ritorno il nuovo id della conversazione convertito in intero
		return int(convID), nil
	}

	return 0, err
}
