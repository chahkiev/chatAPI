package main

import (
	chat "chatAPI/chat"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	// DSN это соединение с базой
	// docker run -p 3306:3306 -v $(PWD):/docker-entrypoint-initdb.d -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=golang -d mysql
	DSN                 = "root:1234@tcp(127.0.0.1:3306)/golang?"
	numberOfFailTryings = 5 // 5 * 1 min delay = 5 minutes
)

func main() {
	db, err := sql.Open("mysql", DSN)
	err = db.Ping() // первое подключение к базе
	tryingsToConnectDB := 0

	for {
		if tryingsToConnectDB == numberOfFailTryings {
			panic("Can't connect to database")
		}
		if err = db.Ping(); err != nil {
			tryingsToConnectDB++
			time.Sleep(time.Minute * 1)
			continue
		}
		break
	}

	chat.InitDB(db)

	chatExplorer := chat.NewChatExplorer(db)
	if err != nil {
		panic(err)
	}

	// s := time.Now().UTC()
	// fmt.Printf("%v \n", s)
	// fmt.Printf("%T \n", s)
	// query := fmt.Sprintf("INSERT INTO %s (username, created_at) VALUES (?, ?)", "users")
	// chatExplorer.DB.Exec(query, "user_6", s)

	r := mux.NewRouter()
	r.HandleFunc("/users/add", chatExplorer.HandlerUserAdd).
		Methods("POST").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/chats/add", chatExplorer.HandlerChatAdd).
		Methods("POST").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/messages/add", chatExplorer.HandlerMessageAdd).
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
