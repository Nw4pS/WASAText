package database

// ConversationPreview modella la riga restituita dal DB per l'elenco chat.
type ConversationPreview struct {
	ID                 int
	Name               string
	Picture            string
	LastMessageDate    string
	LastMessagePreview string
	IsGroup            bool
}

func (db *appdbimpl) GetMyConversations(requesterID string) ([]ConversationPreview, error) {
	query := `
    WITH UltimeConversazioni AS (
        SELECT 
            d.id AS conv_id,
            u_other.username AS conv_name,
            COALESCE(u_other.proPic, '') AS conv_pic,
            0 AS is_group
        FROM Ut_dir my_dir
        JOIN Diretti d ON d.id = my_dir.diretto
        JOIN Ut_dir other_dir ON other_dir.diretto = d.id AND other_dir.utente != my_dir.utente
        JOIN Utenti u_other ON u_other.id = other_dir.utente
        WHERE my_dir.utente = ?

        UNION ALL

        SELECT 
            g.id AS conv_id,
            g.nome AS conv_name,
            COALESCE(g.groupPic, '') AS conv_pic,
            1 AS is_group
        FROM Ut_grup ug
        JOIN Gruppi g ON g.id = ug.gruppo
        WHERE ug.utente = ?
    )
    SELECT 
        c.conv_id,
        c.conv_name,
        c.conv_pic,
        c.is_group,
        COALESCE(m.istante, ''),
        COALESCE(m.contenuto, ''),
        COALESCE(m.is_immagine, 0),
        COALESCE(m.is_eliminato, 0)
    FROM UltimeConversazioni c
    LEFT JOIN Messaggi m ON m.id = (
        SELECT id FROM Messaggi 
        WHERE conversazione = c.conv_id 
        ORDER BY istante DESC, id DESC 
        LIMIT 1
    )
    ORDER BY COALESCE(m.istante, '') DESC;`

	rows, err := db.c.Query(query, requesterID, requesterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Inizializzato a slice vuoto per serializzare in JSON [] invece di null
	convs := make([]ConversationPreview, 0)

	for rows.Next() {
		var c ConversationPreview
		var isGroupInt int
		var lastContent string
		var isImgInt int
		var isDeletedInt int

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Picture,
			&isGroupInt,
			&c.LastMessageDate,
			&lastContent,
			&isImgInt,
			&isDeletedInt,
		)
		if err != nil {
			return nil, err
		}

		c.IsGroup = (isGroupInt == 1)

		// calcolo della preview del messaggio conforme a gocritic
		switch {
		case c.LastMessageDate == "":
			// nessun messaggio scambiato ancora
			c.LastMessagePreview = ""
		case isDeletedInt == 1:
			c.LastMessagePreview = "Questo messaggio è stato eliminato"
		case isImgInt == 1:
			c.LastMessagePreview = "📷 Foto"
		default:
			// se troppo lungo si tronca
			if len(lastContent) > 40 {
				c.LastMessagePreview = lastContent[:40] + "..."
			} else {
				c.LastMessagePreview = lastContent
			}
		}

		convs = append(convs, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return convs, nil
}
