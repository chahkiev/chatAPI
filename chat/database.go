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
	ID    int64   `json:"chat,omitempty"`
	Name  string  `json:"name,omitempty"`
	Users []int64 `json:"users,omitempty"`
}

type ChatUser struct {
	chat int64 `json:"chat,omitempty"`
	user int64 `json:"user,omitempty"`
}

type Message struct {
	Chat   int64  `json:"chat"`
	Author int64  `json:"author"`
	Text   string `json:"text"`
}

func InitDB(db *sql.DB) {
	qs := []string{
		`DROP TABLE IF EXISTS messages;`,
		`DROP TABLE IF EXISTS chat_user;`,
		`DROP TABLE IF EXISTS users;`,
		`DROP TABLE IF EXISTS chats;`,

		`CREATE TABLE users (
id int(11) NOT NULL AUTO_INCREMENT,
username varchar(255) NOT NULL,
created_at datetime NOT NULL,
PRIMARY KEY (id),
CONSTRAINT uUser UNIQUE NONCLUSTERED (
	username
)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE chats (
id int(11) NOT NULL AUTO_INCREMENT,
name varchar(255) NOT NULL,
created_at datetime NOT NULL,
PRIMARY KEY (id),
CONSTRAINT uChat UNIQUE NONCLUSTERED (
	name
)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE chat_user (
id int(11) NOT NULL AUTO_INCREMENT,
chat int(11) NOT NULL,
user int(11) NOT NULL,
PRIMARY KEY (id),
FOREIGN KEY (chat) REFERENCES chats(id),
FOREIGN KEY (user) REFERENCES users(id),
CONSTRAINT uChatUser UNIQUE NONCLUSTERED (
	user, chat
)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE messages (
id int(11) NOT NULL AUTO_INCREMENT,
chat int(11) NOT NULL,
author int(11) NOT NULL,
text varchar(255) NOT NULL,
created_at datetime NOT NULL,
PRIMARY KEY (id),
FOREIGN KEY (chat) REFERENCES chats(id),
FOREIGN KEY (author) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		// INSERTING TEST DATA
		`INSERT IGNORE INTO users (id, username, created_at) VALUES
(1, 'user_1',	'2011-12-18 13:17:17'),
(2, 'user_2',	'2012-12-18 13:17:17'),
(3, 'user_3',	'2011-12-18 13:17:17');`,

		`INSERT IGNORE INTO chats (id, name, created_at) VALUES
(1, 'chat_1', '2012-12-18 13:17:17'),
(2, 'chat_2', '2012-12-18 13:17:17');`,

		`INSERT IGNORE INTO chat_user (chat, user) VALUES
(1, 2),
(1, 3),
(2, 1);`,

		`INSERT INTO messages (id ,chat, author, text, created_at) VALUES
(1, 1, 1, 'some text message', '2012-12-18 13:17:17');`,
	}

	for _, q := range qs {
		_, err := db.Exec(q)
		if err != nil {
			panic(err)
		}
	}
}
