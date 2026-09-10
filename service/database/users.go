package database

// crea utente
func (db *appdbimpl) CreateUser(username string) error {
	_, err := db.c.Exec("INSERT INTO utente (username) VALUES (?)", username)
	return err
}

// aggiorna foto profilo
func (db *appdbimpl) UpdatePropic(username string, proPic string) error {
	_, err := db.c.Exec("UPDATE utenti SET proPic=? WHERE username=?", proPic, username)
	return err
}

// lista gli username di tutti gli utenti di wasatext
func (db *appdbimpl) ListUsernames() ([]string, error) {
	var usernameList []string
	sqlRows, err := db.c.Query("SELECT username FROM utenti")
	for sqlRows.Next() {
		var username string
		if err := sqlRows.Scan(&username); err != nil {
			return nil, err
		}
		usernameList = append(usernameList, username)
	}
	return usernameList, err
}
