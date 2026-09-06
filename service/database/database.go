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

	// Start Database
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

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil { //se l'oggetto che dovrebbe rappresentare il database è vuoto:
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// se il database non è vuvoto: Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) { // il database è vuoto
		sqlStmt := `CREATE TABLE IF NOT EXISTS utenti (
							username TEXT PRIMARY KEY,
							proPic TEXT);
					CREATE TABLE IF NOT EXISTS conversazioni (
							id INTEGER PRIMARY KEY AUTOINCREMENT);
					CREATE TABLE IF NOT EXISTS diretti(
							id INTEGER PRIMARY KEY,
							FOREIGN KEY (id) REFERENCES conversazione(id) ON DELETE CASCADE);
					CREATE TABLE IF NOT EXISTS gruppi(
							nome TEXT NOT NULL,
							groupPic TEXT,
							id INTEGER PRIMARY KEY,
							FOREIGN KEY (id) REFERENCES conversazione(id) ON DELETE CASCADE);
					CREATE TABLE IF NOT EXISTS messaggi(
							utente TEXT NOT NULL,
							conversazione INTEGER NOT NULL,
							istante TEXT NOT NULL,
							contenuto TEXT NOT NULL,
							is_immagine INT NOT NULL,
							is_inoltrato INT NOT NULL,
							is_eliminato INT NOT NULL,
							risposta_a INT,
							stato TEXT CHECK( stato IN ('inviato','ricevuto','letto')) NOT NULL DEFAULT 'inviato',
							id INT NOT NULL,
							FOREIGN KEY(utente) REFERENCES utenti(username),
							FOREIGN KEY(conversazione) REFERENCES conversazioni(id),
							FOREIGN KEY (risposta_a) REFERENCES messaggi(id),
							);
					CREATE TABLE IF NOT EXISTS reaction(
							utente TEXT NOT NULL),
							messaggio INTEGER NOT NULL,
							contenuto TEXT NOT NULL CHECK(length(emoji) >= 1 AND length(emoji) <= 8),
							id INT NOT NULL,
							FOREIGN KEY(utente) REFERENCES utenti(username),
							FOREIGN KEY(messaggio) REFERENCES messaggi(id),
							);
					CREATE TABLE IF NOT EXISTS ut_grup (
							utente TEXT NOT NULL,
							gruppo TEXT NOT NULL,
							FOREIGN KEY (utente) REFERENCES utenti(username),
							FOREIGN KEY (gruppo) REFERENCES gruppi(id));
					CREATE TABLE IF NOT EXISTS ut_dir(
							utente TEXT NOT NULL,
							diretto INTEGER NOT NULL,
							FOREIGN KEY (utente) REFERENCES utenti(username),
							FOREIGN KEY (diretto) REFERENCES diretti(id))
					)`
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
