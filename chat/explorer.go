package chat

import (
	"database/sql"
	"fmt"
	"time"
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
	created_at := time.Now().UTC()
	query := fmt.Sprintf("INSERT INTO users (username, created_at) VALUES (?, ?)")
	result, err := expl.DB.Exec(query, username, created_at)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (expl *ChatExplorer) addChat(chatName string, usersI interface{}) (int64, error) {
	created_at := time.Now().UTC()
	queryForChats := fmt.Sprintf("INSERT INTO chats (name, created_at) VALUES (?, ?)")
	result, err := expl.DB.Exec(queryForChats, chatName, created_at)
	if err != nil {
		return 0, err
	}
	chatId, _ := result.LastInsertId()

	valuesChatUser := []interface{}{}
	users := usersI.([]int64)
	queryRaw := "INSERT INTO chat_user (chat, user) VALUES "

	for _, user := range users {
		queryRaw = queryRaw + " (?, ?),"
		valuesChatUser = append(valuesChatUser, chatId, user)
	}

	runes := []rune(queryRaw)
	query := string(runes[:(len(queryRaw) - 1)])
	query = query + ";"

	expl.DB.Exec(query, valuesChatUser...)

	return chatId, nil
}

func (expl *ChatExplorer) addMessage(chat int64, author int64, text string) (int64, error) {
	created_at := time.Now().UTC()
	query := fmt.Sprintf("INSERT INTO messages (chat, author, text, created_at) VALUES (?, ?, ?, ?)")
	result, err := expl.DB.Exec(query, chat, author, text, created_at)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (expl *ChatExplorer) getChats(user int64) ([]interface{}, error) {
	return []interface{}{1, 2, 3, 4}, nil //it's test data
}

func (expl *ChatExplorer) getMessages(chat int64) ([]interface{}, error) {
	return []interface{}{1, 2, 3, 4}, nil //it's test data
}
