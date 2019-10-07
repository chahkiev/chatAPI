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
	// getting Users Chats
	query := fmt.Sprintf("SELECT chat FROM %s WHERE user = ?", "chat_user")
	rows, err := expl.DB.Query(query, user)
	if err != nil {
		return []interface{}{}, err
	}

	var userChats = []int64{}
	var userMessages = []interface{}{}

	for rows.Next() {
		var chatId int64
		if err := rows.Scan(&chatId); err != nil {
			return []interface{}{}, err
		}
		userChats = append(userChats, chatId)
	}

	// getting last Message from Users Chats
	for _, chat := range userChats {
		query := fmt.Sprintf("SELECT text, created_at FROM %s WHERE chat = ? order by created_at desc limit 1", "messages")
		rows, err := expl.DB.Query(query, chat)
		if err != nil {
			return []interface{}{}, err
		}

		for rows.Next() {
			var message = Message{}
			if err := rows.Scan(&message.Text, &message.Created_at); err != nil {
				return []interface{}{}, err
			}
			message.Chat = chat
			message.Author = user
			userMessages = append(userMessages, message)
		}
	}
	return userMessages, nil
}

func (expl *ChatExplorer) getMessages(chat int64) ([]interface{}, error) {
	return []interface{}{1, 2, 3, 4}, nil //it's test data
}
