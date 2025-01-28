// Функция для загрузки HTML-контента в блок
function loadContent(templatePath) {
    fetch(templatePath)
        .then(response => {
            if (!response.ok) {
                throw new Error("Ошибка загрузки файла: " + response.status);
            }
            return response.text();
        })
        .then(html => {
            document.getElementById("content").innerHTML = html;
        })
        .catch(error => {
            console.error("Ошибка:", error);
            document.getElementById("content").innerHTML = "<p>Ошибка загрузки контента.</p>";
        });
}

// Обработчики кнопок для смены содержимого
document.querySelectorAll('.navbar-button').forEach(button => {
    button.addEventListener('click', () => {
        const page = button.textContent.trim().toLowerCase(); // Имя файла на основе кнопки
        loadContent(`/../templ/${page}.html`);
    });
});

// Загрузка начальной страницы
loadContent("/../templ/statistics.html");