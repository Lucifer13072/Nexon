package main

import (
	"Nexon/admin/scripts"
	a "Nexon/components/sqripts"
	"io/ioutil"
	"log"
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
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("templates/setup/assets"))))

	// Проверка первой настройки
	mux.HandleFunc("/setup", scripts.setupCompleted)

	// Главная страница
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Если настройка не завершена, перенаправляем на страницу настройки
		if !scripts.setupCompleted {
			http.Redirect(w, r, "/setup", http.StatusSeeOther)
			return
		}

		// Пример контента с тегами, который нужно обработать
		pageContent, err := ioutil.ReadFile("templates/main/index.html") // Путь к файлу с тегами
		if err != nil {
			http.Error(w, "Не удалось загрузить страницу", http.StatusInternalServerError)
			log.Printf("Ошибка чтения файла: %v", err)
			return
		}

		// Обрабатываем контент с заменой тегов на реальные данные
		pageContentStr := string(pageContent)
		finalContent := a.RenderContent(pageContentStr)
		// Отправляем отрендеренный контент пользователю
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(finalContent))
	})

	// Запуск сервера
	http.ListenAndServe(":"+port, mux)
}
