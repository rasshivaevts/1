# Учёт личных трат — веб-приложение на Go

Учебный проект по веб-разработке на Go. Этапы 1–2: HTTP-сервер, шаблоны, список трат.

## Технологии
- Go (стандартная библиотека, без сторонних фреймворков)
- net/http — HTTP-сервер и маршрутизация
- html/template — серверный рендеринг HTML
- Веб-интерфейс: HTML + CSS

## Структура проекта

```
expenses_tracker/
├── go.mod                          # Модуль Go
├── main.go                         # Точка входа, запуск сервера
├── .gitignore                      # Исключения для git
├── internal/
│   ├── handlers/
│   │   └── handlers.go             # Все HTTP-обработчики
│   └── models/
│       └── expense.go              # Модель Expense + хранилище в памяти
└── web/
    ├── static/
    │   └── style.css               # Стили интерфейса
    └── templates/
        └── layout.html             # HTML-шаблон (layout + все страницы)
```

## Запуск

```bash
cd expenses_tracker
go run .
```

Сервер поднимется на http://localhost:8080

## Маршруты

| Путь            | Метод | Описание                          |
|-----------------|-------|-----------------------------------|
| `/`             | GET   | Главная страница                  |
| `/about`        | GET   | Описание проекта                  |
| `/ping`         | GET   | Health-check → `pong` (иначе 405) |
| `/expenses`     | GET   | Список трат                       |
| `/expenses/new` | GET   | Форма добавления траты            |
| `/expenses`     | POST  | Создание траты → redirect на список|

## Загрузка на GitHub

```bash
git init
git add .
git commit -m "Этапы 1-2: HTTP-сервер, шаблоны, список трат"
git remote add origin https://github.com/USERNAME/expenses_tracker.git
git branch -M main
git push -u origin main
```

Замените `USERNAME` на свой логин GitHub.
