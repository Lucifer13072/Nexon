package sqripts

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// Моковые функции для получения данных, их нужно будет заменить на реальные запросы к базе данных
func getNews() string {
	// Замени на реальную логику получения новостей из базы данных
	return "Последние новости на сайте"
}

func getLoginForm() string {
	fileContent, err := ioutil.ReadFile("templates/forms/login_form.html") // Заменить путь на актуальный
	if err != nil {
		log.Println("Ошибка при чтении HTML файла:", err)
		return fmt.Sprintf("<p>Ошибка загрузки формы регистрации</p>")
	}

	// Возвращаем содержимое HTML файла как строку
	return string(fileContent)
}

func getRegistrationForm() string {
	fileContent, err := ioutil.ReadFile("templates/forms/registration_form.html") // Заменить путь на актуальный
	if err != nil {
		log.Println("Ошибка при чтении HTML файла:", err)
		return fmt.Sprintf("<p>Ошибка загрузки формы регистрации</p>")
	}

	// Возвращаем содержимое HTML файла как строку
	return string(fileContent)
}

func getContent() string {
	// Возвращаем динамический контент
	return "Динамический контент для авторизованных пользователей"
}

// Функция проверки авторизации (используй свою логику проверки сессий или куки)
func isUserLoggedIn() bool {
	// В реальном коде здесь будет проверка сессии или куки
	return true // Для примера всегда авторизован
}

// Функция для получения группы пользователя
func getUserGroup() int {
	// В реальном коде здесь будет запрос к базе данных для получения группы пользователя
	return 1 // Для примера возвращаем группу 1
}

// Функция для рендеринга контента с заменой тегов
func RenderContent(content string) string {
	// Заменяем стандартные теги на данные
	content = strings.Replace(content, "{news}", getNews(), -1)
	content = strings.Replace(content, "{login}", getLoginForm(), -1)
	content = strings.Replace(content, "{registration}", getRegistrationForm(), -1)
	content = strings.Replace(content, "{content}", getContent(), -1)

	// Проверка авторизации для тега [login]
	if isUserLoggedIn() {
		content = strings.Replace(content, "[login]", getContent(), -1)
	} else {
		content = strings.Replace(content, "[login]", "", -1)
	}

	// Проверка на группу для тега [login=group]
	group := getUserGroup()
	content = strings.Replace(content, "[login="+string(group)+"]", getContent(), -1)

	return content
}
