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

	// Проверка первой настройки
	mux.HandleFunc("/setup", setupHandler)

	// Главная страница
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !setupCompleted {
			http.Redirect(w, r, "templates/pages/setup", http.StatusSeeOther)
			return
		}
		w.Write([]byte("<h1>Добро пожаловать!</h1>"))
	})
	При этом давай еще файл setup.html и styles_setup.css будут удаляться после установки
	http.ListenAndServe(":"+port, mux)
}
