/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	//+ Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB, qui vengono specificati i metodi che il database dovrà avere
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error
	DoLogin(username string) (string, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil { // se l'oggetto che dovrebbe rappresentare il database è vuoto:
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// se il database non è vuoto: Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) { // il database è vuoto
		sqlStmt := `CREATE TABLE IF NOT EXISTS utenti (
							id TEXT PRIMARY KEY,
							username TEXT UNIQUE,
							proPic TEXT
							);
					CREATE TABLE IF NOT EXISTS conversazioni (
							id INTEGER PRIMARY KEY AUTOINCREMENT
							);
					CREATE TABLE IF NOT EXISTS diretti(
							id INTEGER PRIMARY KEY,
							FOREIGN KEY (id) REFERENCES conversazioni(id) ON DELETE CASCADE
							);
					CREATE TABLE IF NOT EXISTS gruppi(
							nome TEXT NOT NULL,
							groupPic TEXT,
							id INTEGER PRIMARY KEY,
							FOREIGN KEY (id) REFERENCES conversazioni(id) ON DELETE CASCADE
							);
					CREATE TABLE IF NOT EXISTS messaggi(
							utente TEXT NOT NULL,
							conversazione INTEGER NOT NULL,
							istante TEXT NOT NULL,
							contenuto TEXT NOT NULL,
							is_immagine INTEGER NOT NULL,
							is_inoltrato INTEGER NOT NULL,
							is_eliminato INTEGER NOT NULL,
							risposta_a INTEGER,
							stato TEXT CHECK( stato IN ('sent','received','read')) NOT NULL DEFAULT 'sent',
							id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
							FOREIGN KEY(utente) REFERENCES utenti(id),
							FOREIGN KEY(conversazione) REFERENCES conversazioni(id),
							FOREIGN KEY (risposta_a) REFERENCES messaggi(id)
							);
					CREATE TABLE IF NOT EXISTS reaction(
							utente TEXT NOT NULL,
							messaggio INTEGER NOT NULL,
							contenuto TEXT NOT NULL CHECK(length(contenuto) >= 1 AND length(contenuto) <= 8),
							id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
							FOREIGN KEY(utente) REFERENCES utenti(id),
							FOREIGN KEY(messaggio) REFERENCES messaggi(id)
							);
					CREATE TABLE IF NOT EXISTS ut_grup (
							utente TEXT NOT NULL,
							gruppo INTEGER NOT NULL,
							PRIMARY KEY(utente,gruppo),
							FOREIGN KEY (utente) REFERENCES utenti(id),
							FOREIGN KEY (gruppo) REFERENCES gruppi(id)
							);
					CREATE TABLE IF NOT EXISTS ut_dir(
							utente TEXT NOT NULL,
							diretto INTEGER NOT NULL,
							PRIMARY KEY(utente,diretto),
							FOREIGN KEY (utente) REFERENCES utenti(id),
							FOREIGN KEY (diretto) REFERENCES diretti(id)
							);
					
					CREATE TRIGGER IF NOT EXISTS trg_prevent_duplicate_id_group
						BEFORE INSERT ON gruppi
						FOR EACH ROW
						WHEN EXISTS (SELECT 1 FROM diretti WHERE id = NEW.id)
						BEGIN
						    SELECT RAISE(ABORT, 'Questo ID conversazione è già utilizzato per una chat diretta.');
						END;
					
					CREATE TRIGGER IF NOT EXISTS trg_prevent_duplicate_id_direct
						BEFORE INSERT ON diretti
						FOR EACH ROW
						WHEN EXISTS (SELECT 1 FROM gruppi WHERE id = NEW.id)
						BEGIN
						    SELECT RAISE(ABORT, 'Questo ID conversazione è già utilizzato per un gruppo.');
						END;

					CREATE TRIGGER IF NOT EXISTS trg_limit_ut_dir_insert
						BEFORE INSERT ON ut_dir
						FOR EACH ROW
						WHEN (SELECT COUNT(*) FROM ut_dir WHERE diretto = NEW.diretto) >= 2
						BEGIN
						    SELECT RAISE(ABORT, 'Una conversazione diretta non può avere più di due utenti.');
						END;
					CREATE TRIGGER IF NOT EXISTS trg_prevent_ut_dir_delete
						BEFORE DELETE ON ut_dir
						FOR EACH ROW
						BEGIN
							SELECT RAISE(ABORT, 'Non è permesso rimuovere utenti da una conversazione diretta.');
						END;
					
					CREATE TRIGGER IF NOT EXISTS trg_prevent_empty_group
						BEFORE DELETE ON ut_grup
						FOR EACH ROW
						WHEN (SELECT COUNT(*) FROM ut_grup WHERE gruppo = OLD.gruppo) <= 1
						BEGIN
							SELECT RAISE(ABORT, 'Un gruppo deve avere almeno un utente associato.');
						END;
					
					CREATE TRIGGER IF NOT EXISTS trg_cascade_delete_gruppo_to_conv
						AFTER DELETE ON gruppi
						FOR EACH ROW
						BEGIN
						    DELETE FROM conversazioni WHERE id = OLD.id;
						END;
					
					CREATE TRIGGER IF NOT EXISTS trg_cascade_delete_diretto_to_conv
						AFTER DELETE ON diretti
						FOR EACH ROW
						BEGIN
							DELETE FROM conversazioni WHERE id = OLD.id;
						END;
					`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
