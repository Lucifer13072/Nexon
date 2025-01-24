package main

import (
	"Nexon/admin/adminScripts"
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
	mux.Handle("/configs/", http.StripPrefix("/configs", http.FileServer(http.Dir("admin/adminScripts/configs"))))
	// Проверка первой настройки
	mux.HandleFunc("/setup", adminScripts.SetupHandler)

	// Главная страница
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Если настройка не завершена, перенаправляем на страницу настройки
		if !adminScripts.SetupComleteRead() {
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
		finalContent := pageContentStr
		// Отправляем отрендеренный контент пользователю
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(finalContent))
	})

	// Авторизация в админ панель
	mux.HandleFunc("/admin/login", adminScripts.AdminLoginHandler)

	// Админ панель (доступ только для авторизованных пользователей)
	mux.HandleFunc("/admin/index", func(w http.ResponseWriter, r *http.Request) {
		// Проверка авторизации с использованием куки
		if !adminScripts.IsUserAuthenticated(r) {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		// Загружаем страницу админ панели
		pageContent, err := ioutil.ReadFile("admin/front/dashboard.html")
		if err != nil {
			http.Error(w, "Не удалось загрузить страницу админ панели", http.StatusInternalServerError)
			log.Printf("Ошибка чтения файла: %v", err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(pageContent)
	})

	// Запуск сервера
	http.ListenAndServe(":"+port, mux)
}
