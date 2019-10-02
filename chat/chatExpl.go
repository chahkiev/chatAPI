package chat

import (
	"database/sql"
	"net/http"
)

type ChatExplorer struct {
	DB *sql.DB
}

func NewChatExplorer(db *sql.DB) *ChatExplorer {
	return &ChatExplorer{
		DB: db,
	}
}

func (chatExpl *ChatExplorer) HandlerUsersAdd(w http.ResponseWriter, r *http.Request) {
}

func (chatExpl *ChatExplorer) HandlerChatssAdd(w http.ResponseWriter, r *http.Request) {
}

func (chatExpl *ChatExplorer) HandlerMessagesAdd(w http.ResponseWriter, r *http.Request) {
}

func (chatExpl *ChatExplorer) HandlerChatsGet(w http.ResponseWriter, r *http.Request) {
}

func (chatExpl *ChatExplorer) HandlerMessagesGet(w http.ResponseWriter, r *http.Request) {
}
