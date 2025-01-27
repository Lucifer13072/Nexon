package adminScripts

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Проверка аутентификации
func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Проверка куки на авторизацию
		cookie, err := r.Cookie("admin_auth")
		if err == nil {
			if ValidateCookie(cookie.Value) {
				http.Redirect(w, r, "/admin/index", http.StatusSeeOther)
				return
			}
		}

		tmpl, err := template.ParseFiles("admin/front/login.html")
		if err != nil {
			http.Error(w, "Не удалось загрузить шаблон", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Проверка логина и пароля
		if AuthenticateAdmin(username, password) {
			// Успешная авторизация, установка куки
			token := GenerateAuthToken(username)

			// Сохраняем токен в базе данных
			err := SaveTokenInSession(username, token)
			if err != nil {
				http.Error(w, "Ошибка при сохранении токена", http.StatusInternalServerError)
				return
			}

			// Устанавливаем куки
			cookie := http.Cookie{
				Name:     "admin_auth",
				Value:    token,
				Expires:  time.Now().Add(24 * time.Hour),
				Path:     "/",
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)

			// Перенаправление в админ панель
			http.Redirect(w, r, "/admin/index", http.StatusSeeOther)
			return
		} else {
			// Ошибка авторизации
			tmpl, err := template.ParseFiles("admin/front/login.html")
			if err != nil {
				http.Error(w, "Не удалось загрузить шаблон", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, map[string]string{"Error": "Неверный логин или пароль"})
			return
		}
	}
}

// Функция для проверки логина и пароля
func AuthenticateAdmin(username, password string) bool {
	dbIp, dbUser, dbPassword, dbName := DatabaseSettingsReader()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbIp, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Ошибка подключения к базе данных:", err)
		return false
	}
	defer db.Close()

	passwordHash := HashPassword(password)

	var groupUser int
	err = db.QueryRow("SELECT group_user FROM users WHERE username = ? AND password_hash = ?", username, passwordHash).Scan(&groupUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Println("Ошибка выполнения запроса:", err)
		return false
	}

	return groupUser == 0 || groupUser == 1
}

// Генерация токена для куки
func GenerateAuthToken(username string) string {
	// Простая генерация токена (например, через хеширование)
	return HashPassword(username + fmt.Sprint(time.Now().Unix()))
}

// Проверка валидности токена из куки
func ValidateCookie(token string) bool {
	// Чтение настроек базы данных
	dbIp, dbUser, dbPassword, dbName := DatabaseSettingsReader()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbIp, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Ошибка подключения к базе данных:", err)
		return false
	}
	defer db.Close()

	// Проверяем, существует ли токен в таблице user_sessions
	var storedToken string
	var expiration time.Time
	err = db.QueryRow("SELECT token, expiration FROM user_sessions WHERE token = ?", token).Scan(&storedToken, &expiration)
	if err != nil {
		if err == sql.ErrNoRows {
			// Токен не найден
			return false
		}
		log.Println("Ошибка выполнения запроса:", err)
		return false
	}

	// Проверяем, не истек ли срок действия токена
	if time.Now().After(expiration) {
		// Токен истек
		return false
	}

	// Токен найден и не истек
	return true
}

func IsUserAuthenticated(r *http.Request) bool {
	// Проверка наличия куки
	cookie, err := r.Cookie("admin_auth")
	if err != nil {
		return false
	}

	// Проверка валидности куки
	return ValidateCookie(cookie.Value)
}

func SaveTokenInSession(username, token string) error {
	// Чтение настроек базы данных
	dbIp, dbUser, dbPassword, dbName := DatabaseSettingsReader()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbIp, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Ошибка подключения к базе данных:", err)
		return err
	}
	defer db.Close()

	// Устанавливаем срок действия токена (например, 24 часа)
	expiration := time.Now().Add(24 * time.Hour)

	// Вставляем или обновляем запись с токеном в базе данных
	_, err = db.Exec("INSERT INTO user_sessions (username, token, expiration) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE token = ?, expiration = ?",
		username, token, expiration, token, expiration)

	if err != nil {
		log.Println("Ошибка сохранения токена в базу данных:", err)
		return err
	}

	return nil
}
