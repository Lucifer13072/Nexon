package adminScripts

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"time"
)

// Данные для подключения к базе данных
var (
	dbIP, dbUser, dbPassword, dbBase = DatabaseSettingsReader()
)

// Обработчик для логина
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Чтение данных из формы
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Подключение к базе данных
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbIP + ":3306)/" + dbBase
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Проверка имени пользователя и пароля
	var passwordHash string
	var group int
	err = db.QueryRow("SELECT password_hash, group FROM users WHERE username = ?", username).Scan(&passwordHash, &group)
	if err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	// Проверка хэша пароля
	if HashPassword(password) != passwordHash {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	// Проверка группы пользователя
	if group != 0 && group != 1 {
		http.Error(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	// Создание куки
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    "secure-token-value", // Здесь можно использовать JWT или другой способ генерации токена
		Expires:  expiration,
		HttpOnly: true,
		Secure:   true, // Используйте true, если работаете через HTTPS
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Успешная авторизация"))
}

// Middleware для проверки авторизации
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie.Value != "secure-token-value" {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Пример защищенного маршрута
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Доступ к защищенному ресурсу"))
}
