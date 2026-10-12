package database

import (
	"database/sql"
	"fmt"
)

// SetPropic (o SetMyPhoto) aggiorna la foto profilo
func (db *appdbimpl) SetPropic(identifier string, proPic string) error {
	res, err := db.c.Exec("UPDATE Utenti SET proPic = ? WHERE id = ?", proPic, identifier)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiornamento della foto: %w", err)
	}

	// salvo quante tuple sono state modificate, per poter sapere se l'update ha avuto successo
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		//nessuna tupla è stata modificata
		return sql.ErrNoRows
	}

	return nil
}
