package chat

import (
	"database/sql"
)

type ChatExplorer struct {
	DB *sql.DB
}

type response struct {
	Response interface{} `json:"response,omitempty"`
}

type errorResp struct {
	Error string `json:"error,omitempty"`
}

func NewChatExplorer(db *sql.DB) *ChatExplorer {
	return &ChatExplorer{
		DB: db,
	}
}

func (expl *ChatExplorer) addUser(username string) (int64, error) {
	return 1, nil //it's test data
}

func (expl *ChatExplorer) addChat(chatName string, users interface{}) (int64, error) {
	return 1, nil //it's test data
}

func (expl *ChatExplorer) addMessage(chat int64, author int64, text string) (int64, error) {
	return 1, nil //it's test data
}

func (expl *ChatExplorer) getChats(user int64) ([]interface{}, error) {
	return []interface{}{1, 2, 3, 4}, nil //it's test data
}

func (expl *ChatExplorer) getMessages(chat int64) ([]interface{}, error) {
	return []interface{}{1, 2, 3, 4}, nil //it's test data
}
