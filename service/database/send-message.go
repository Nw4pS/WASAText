package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type SendMessageRequest struct {
	ConversationID int
	SenderID       string
	Content        string
	IsImage        bool
	ReplyTo        *int // uso un puntatore per gestire i valori NULL del DB
	IsForwarded    bool // true se è un messaggio inoltrato
}

func (db *appdbimpl) SendMessage(req SendMessageRequest) (int, error) {
	// verifico se l'utente sta inviando un messaggio ad un gruppo o ad una chat diretta, in ogni caso controllo se l'utente è nella chat
	var isMember int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Ut_dir WHERE utente = ? AND diretto = ?", req.SenderID, req.ConversationID).Scan(&isMember)
	if err != nil {
		return 0, fmt.Errorf("errore durante la verifica nella chat diretta: %w", err)
	}

	if isMember == 0 {
		err = db.c.QueryRow("SELECT COUNT(*) FROM Ut_grup WHERE utente = ? AND gruppo = ?", req.SenderID, req.ConversationID).Scan(&isMember)
		if err != nil {
			return 0, fmt.Errorf("errore durante la verifica nel gruppo: %w", err)
		}
	}

	// se l'utente non è membro di nessuna delle due, nego l'accesso
	if isMember == 0 {
		return 0, errors.New("operazione negata: l'utente non fa parte della conversazione")
	}

	// inserimento del messaggio
	currentTime := time.Now().UTC().Format(time.RFC3339)
	var res sql.Result

	if req.ReplyTo != nil {
		res, err = db.c.Exec(
			"INSERT INTO Messaggi (utente, conversazione, istante, contenuto, is_immagine, is_inoltrato, is_eliminato, risposta_a, stato) VALUES (?, ?, ?, ?, ?, 0, 0, ?, 'sent')",
			req.SenderID, req.ConversationID, currentTime, req.Content, req.IsImage, *req.ReplyTo,
		)
	} else {
		res, err = db.c.Exec(
			"INSERT INTO Messaggi (utente, conversazione, istante, contenuto, is_immagine, is_inoltrato, is_eliminato, stato) VALUES (?, ?, ?, ?, ?, 0, 0, 'sent')",
			req.SenderID, req.ConversationID, currentTime, req.Content, req.IsImage,
		)
	}

	if err != nil {
		return 0, fmt.Errorf("errore durante l'inserimento del messaggio: %w", err)
	}

	messageID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("errore nel recupero dell'ID del messaggio: %w", err)
	}

	return int(messageID), nil
}
