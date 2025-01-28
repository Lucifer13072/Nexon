document.addEventListener("DOMContentLoaded", function () {
    const ctx = document.getElementById('visitsChart').getContext('2d');

    // Данные для графика
    const visitsData = {
        labels: ['Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница', 'Суббота', 'Воскресенье'],
        datasets: [{
            label: 'Посещения сайта',
            data: [120, 200, 150, 300, 400, 500, 250], // Примерные данные
            backgroundColor: 'rgba(54, 162, 235, 0.2)',
            borderColor: 'rgba(54, 162, 235, 1)',
            borderWidth: 1
        }]
    };

    // Конфигурация графика
    const config = {
        type: 'line', // Линейный график
        data: visitsData,
        options: {
            responsive: true,
            plugins: {
                legend: {
                    display: true,
                    position: 'top'
                }
            },
            scales: {
                y: {
                    beginAtZero: true // Начало оси Y с нуля
                }
            }
        }
    };

    // Создание графика
    new Chart(ctx, config);
});