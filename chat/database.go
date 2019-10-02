package chat

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

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
