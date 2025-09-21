package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Todo struct {
	ID       int
	Title    string
	Done     bool
	Deadline string
}

var db *sql.DB

func InitDB() {
	var err error
	connStr := "host=db port=5432 user=todo password=password dbname=todo sslmode=disable"

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil && db.Ping() == nil {
			fmt.Println("✅ DB接続成功")
			return
		}
		fmt.Println("DB起動待機中...再試行")
		time.Sleep(2 * time.Second)
	}
	log.Fatal("DB接続失敗: ", err)
}

func GetTodos() ([]Todo, error) {
	rows, err := db.Query("SELECT id, title, done, COALESCE(deadline::text, '') FROM todos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.Deadline); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

// 新規Todo登録
func InsertTodo(title, deadline string) error {
	_, err := db.Exec("INSERT INTO todos (title, deadline) VALUES ($1, $2)", title, deadline)
	return err
}

func DeleteTodo(id int) error {
	_, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
	return err
}
