package chat

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int64  `json:"user,omitempty"`
	Username string `json:"username,omitempty"`
}

type Chat struct {
	ID    int64  `json:"chat,omitempty"`
	Name  string `json:"name,omitempty"`
	Users []int  `json:"users,omitempty"`
}

type Message struct {
	Chat   int64  `json:"chat"`
	Author int64  `json:"author"`
	Text   string `json:"text"`
}

func InitDB(db *sql.DB) {
	qs := []string{
		`DROP TABLE IF EXISTS users;`,

		`CREATE TABLE users (
// TODO
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`DROP TABLE IF EXISTS chats;`,

		`CREATE TABLE chats (
// TODO
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`DROP TABLE IF EXISTS messages;`,

		`CREATE TABLE messages (
// TODO
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	}

	for _, q := range qs {
		_, err := db.Exec(q)
		if err != nil {
			panic(err)
		}
	}
}
