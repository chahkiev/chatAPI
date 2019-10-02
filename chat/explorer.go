package chat

import (
	"database/sql"
)

type ChatExplorer struct {
	DB *sql.DB
}

func NewChatExplorer(db *sql.DB) *ChatExplorer {
	return &ChatExplorer{
		DB: db,
	}
}
