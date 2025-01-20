package main

import (
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()

	// Обработчик для статики (например, CSS, изображения и т.д.)
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("templates/pages/"))))

	// Проверка первой настройки
	mux.HandleFunc("/setup", setupHandler)

	// Главная страница
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !setupCompleted {
			http.Redirect(w, r, "/setup", http.StatusSeeOther)
			return
		}
		w.Write([]byte("<h1>Добро пожаловать!</h1>"))
	})

	// Запуск сервера
	http.ListenAndServe(":"+port, mux)
}
