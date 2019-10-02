package main

import (
	chat "chatAPI/chat"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	// DSN это соединение с базой
	// docker run -p 3306:3306 -v $(PWD):/docker-entrypoint-initdb.d -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=golang -d mysql
	DSN = "root:1234@tcp(127.0.0.1:3306)/golang?"
)

func main() {
	db, err := sql.Open("mysql", DSN)
	err = db.Ping() // первое подключение к базе
	if err != nil {
		panic(err)
	}

	chat.InitDB(db)

	chatExplorer := chat.NewChatExplorer(db)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/users/add", chatExplorer.HandlerUsersAdd).
		Methods("POST").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/chats/add", chatExplorer.HandlerChatssAdd).
		Methods("POST").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/messages/add", chatExplorer.HandlerMessagesAdd).
		Methods("POST").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/chats/get", chatExplorer.HandlerChatsGet).
		Methods("POST").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/messages/get", chatExplorer.HandlerMessagesGet).
		Methods("POST").
		Headers("Content-Type", "application/json")

	fmt.Println("starting server at :9000")
	log.Fatal(http.ListenAndServe(":9000", r))
}
