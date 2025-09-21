package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"todoapp/models"
)

type Todo struct {
	ID       int
	Title    string
	Done     bool
	Deadline string
}

var todos = []Todo{
	{Title: "Learn Go", Done: true},
	{Title: "Build Todo App", Done: false},
	{Title: "Deploy", Done: false},
}

// 削除ハンドラー
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォーム解析失敗", http.StatusBadRequest)
		return
	}

	// フォームから ID を取得
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID が不正です", http.StatusBadRequest)
		return
	}

	// DB から削除
	if err := models.DeleteTodo(id); err != nil {
		log.Println("削除エラー:", err)
		http.Error(w, "削除失敗: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// フォームデータをパース
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}
	todo := r.FormValue("todo")
	deadline := r.FormValue("deadline")

	if deadline == nil || todo == nil {
		http.Error(w, "データを入力してください。", http.StatusInternalServerError)
	}

	if err := models.InsertTodo(todo, deadline); err != nil {
		http.Error(w, "DB登録失敗: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	models.InitDB()
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("../templates/index.html"))

		todos, err := models.GetTodos()
		if err != nil {
			http.Error(w, "DB取得エラー: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Title": "Todo list",
			"Todos": todos,
		}

		tmpl.Execute(w, data)
	})

	// React用Cors
	// handler := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"http://localhost:3000"},
	// 	AllowedMethods: []string{http.MethodGet, http.MethodPost, http.Methodput, http.MethodDelete},
	// 	AllowedHeaders: []string{"*"},
	// }).Handler(mux)

	mux.HandleFunc("/submit", submitHandler)
	mux.HandleFunc("/delete", deleteHandler)
	http.ListenAndServe(":8080", mux)
}
