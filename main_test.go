package main

import (
	"bytes"
	chat "chatAPI/chat"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"net/http"
	"net/http/httptest"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// CaseResponse
type CR map[string]interface{}

type Case struct {
	ts     httptest.Server
	Method string // GET по-умолчанию в http.NewRequest если передали пустую строку
	Path   string
	Query  string
	Status int
	Result interface{}
	Body   interface{}
}

var (
	client              = &http.Client{Timeout: time.Second}
	numberOfFailTryings = 5 // 5 * 1 min delay = 5 minutes
)

func PrepareTestApis(db *sql.DB) {
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
(2, 1),
(1, 1);`,

		`INSERT INTO messages (id ,chat, author, text, created_at) VALUES
(1, 1, 1, 'some', '2012-12-18 19:17:17'),
(2, 1, 1, 'text', '2012-12-18 17:17:17'),
(3, 2, 1, 'message', '2012-12-18 13:17:17');`,
	}

	for _, q := range qs {
		_, err := db.Exec(q)
		if err != nil {
			panic(err)
		}
	}
}

func CleanupTestApis(db *sql.DB) {
	qs := []string{
		`DROP TABLE IF EXISTS messages;`,
		`DROP TABLE IF EXISTS chat_user;`,
		`DROP TABLE IF EXISTS users;`,
		`DROP TABLE IF EXISTS chats;`,
	}
	for _, q := range qs {
		_, err := db.Exec(q)
		if err != nil {
			panic(err)

		}
	}
}

func TestApis(t *testing.T) {
	db, err := sql.Open("mysql", DSN)
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

	PrepareTestApis(db)
	defer CleanupTestApis(db)

	handler := chat.NewChatExplorer(db)

	// ts := httptest.NewServer(http.HandlerFunc(handler.HandlerChatAdd))

	cases := []Case{
		Case{
			ts:     *httptest.NewServer(http.HandlerFunc(handler.HandlerUserAdd)),
			Path:   "/users/add", // список таблиц
			Method: http.MethodPost,
			Body: CR{
				"username": "user_4",
			},
			Result: CR{
				"response": CR{
					"id": 4,
				},
			},
		},
		Case{
			ts:     *httptest.NewServer(http.HandlerFunc(handler.HandlerChatAdd)),
			Path:   "/chats/add", // список таблиц
			Method: http.MethodPost,
			Body: CR{
				"name":  "chat_3",
				"users": []int64{1, 2},
			},

			Result: CR{
				"response": CR{
					"id": 3,
				},
			},
		},
		Case{
			ts:     *httptest.NewServer(http.HandlerFunc(handler.HandlerMessageAdd)),
			Path:   "/messages/add", // список таблиц
			Method: http.MethodPost,
			Body: CR{
				"chat":   1,
				"author": 1,
				"test":   "test text",
			},

			Result: CR{
				"response": CR{
					"id": 4,
				},
			},
		},
	}

	runCases(t, db, cases)
}

func runCases(t *testing.T, db *sql.DB, cases []Case) {
	for idx, item := range cases {
		var (
			err      error
			result   interface{}
			expected interface{}
			req      *http.Request
		)

		caseName := fmt.Sprintf("case %d: [%s] %s %s", idx, item.Method, item.Path, item.Query)

		if db.Stats().OpenConnections != 1 {
			t.Fatalf("[%s] you have %d open connections, must be 1", caseName, db.Stats().OpenConnections)
		}

		data, err := json.Marshal(item.Body)
		if err != nil {
			panic(err)
		}
		reqBody := bytes.NewReader(data)
		req, err = http.NewRequest(item.Method, item.ts.URL+item.Path, reqBody)
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("[%s] request error: %v", caseName, err)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if item.Status == 0 {
			item.Status = http.StatusOK
		}

		if resp.StatusCode != item.Status {
			t.Fatalf("[%s] expected http status %v, got %v", caseName, item.Status, resp.StatusCode)
			continue
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Fatalf("[%s] cant unpack json: %v", caseName, err)
			continue
		}

		// reflect.DeepEqual не работает если нам приходят разные типы
		// а там приходят разные типы (string VS interface{}) по сравнению с тем что в ожидаемом результате
		// это конвертит данные сначала в json, а потом обратно в interface - получаем совместимые результаты
		dataR, err := json.Marshal(item.Result)
		json.Unmarshal(dataR, &expected)

		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("[%s] results not match\nGot : %#v\nWant: %#v", caseName, result, expected)
			continue
		}
	}

}
