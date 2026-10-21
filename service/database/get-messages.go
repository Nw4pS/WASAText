package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// Strutture conformi allo schema JSON richiesto dall'API
type Reaction struct {
	ID      int    `json:"id"`
	User    string `json:"user"`
	Content string `json:"content"`
}

type Message struct {
	ID           int        `json:"id"`
	User         string     `json:"user"`
	Conversation int        `json:"conversation"`
	Timestamp    string     `json:"timestamp"`
	Content      string     `json:"content"`
	IsImage      bool       `json:"isImage"`
	IsForwarded  bool       `json:"isForwarded"`
	AnswerTo     *int       `json:"answerTo,omitempty"` // Puntatore per gestire i Null
	IsDeleted    bool       `json:"isDeleted"`
	State        string     `json:"state"`
	Reactions    []Reaction `json:"reactions"`
}

func (db *appdbimpl) GetConversation(requesterID string, convID int) ([]Message, error) {
	// controllo se l'utente fa parte della conversazione
	var isMember int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Ut_dir WHERE utente = ? AND diretto = ?", requesterID, convID).Scan(&isMember)
	if err != nil {
		return nil, fmt.Errorf("errore durante la verifica nella chat diretta: %w", err)
	}

	if isMember == 0 {
		err = db.c.QueryRow("SELECT COUNT(*) FROM Ut_grup WHERE utente = ? AND gruppo = ?", requesterID, convID).Scan(&isMember)
		if err != nil {
			return nil, fmt.Errorf("errore durante la verifica nel gruppo: %w", err)
		}
	}

	if isMember == 0 {
		return nil, errors.New("operazione negata: l'utente non fa parte della conversazione")
	}

	// estraggo i messaggi in ordine cronologico inverso
	// faccio una join con Utenti per restituire lo username richiesto dal frontend[cite: 13, 20]
	rows, err := db.c.Query(`
		SELECT m.id, u.username, m.conversazione, m.istante, m.contenuto, 
		       m.is_immagine, m.is_inoltrato, m.is_eliminato, m.risposta_a, m.stato
		FROM Messaggi m
		JOIN Utenti u ON m.utente = u.id
		WHERE m.conversazione = ?
		ORDER BY m.istante DESC
	`, convID)

	if err != nil {
		return nil, fmt.Errorf("errore nel recupero dei messaggi: %w", err)
	}
	defer rows.Close()

	var messages []Message

	for rows.Next() {
		var msg Message
		var isImg, isFwd, isDel int
		var answerTo sql.NullInt64

		err = rows.Scan(
			&msg.ID, &msg.User, &msg.Conversation, &msg.Timestamp, &msg.Content,
			&isImg, &isFwd, &isDel, &answerTo, &msg.State,
		)
		if err != nil {
			return nil, fmt.Errorf("errore nello scan del messaggio: %w", err)
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("errore durante l'iterazione delle righe: %w", err)
		}

		// Conversione tipi SQLite (INTEGER) in tipi Go (bool / *int)[cite: 20]
		msg.IsImage = isImg != 0
		msg.IsForwarded = isFwd != 0
		msg.IsDeleted = isDel != 0
		if answerTo.Valid {
			val := int(answerTo.Int64)
			msg.AnswerTo = &val
		}

		// 3. Estrazione delle reazioni per questo messaggio
		msg.Reactions, err = db.getReactionsForMessage(msg.ID)
		if err != nil {
			return nil, fmt.Errorf("errore nel recupero delle reactions: %w", err)
		}

		messages = append(messages, msg)
	}

	// Se non ci sono messaggi, inizializza un array vuoto anziché nil per rispettare il JSON "[]"
	if messages == nil {
		messages = []Message{}
	}

	return messages, nil
}

// Funzione helper per recuperare le reazioni con lo username associato[cite: 20]
func (db *appdbimpl) getReactionsForMessage(msgID int) ([]Reaction, error) {
	rows, err := db.c.Query(`
		SELECT r.id, u.username, r.contenuto 
		FROM Reaction r
		JOIN Utenti u ON r.utente = u.id
		WHERE r.Messaggio = ?
	`, msgID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []Reaction
	for rows.Next() {
		var react Reaction
		if err := rows.Scan(&react.ID, &react.User, &react.Content); err != nil {
			return nil, err
		}
		reactions = append(reactions, react)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("errore durante l'iterazione delle righe: %w", err)
	}

	if reactions == nil {
		reactions = []Reaction{}
	}
	return reactions, nil
}
