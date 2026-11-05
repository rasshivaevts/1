package main

import (
 "database/sql"
 "html/template"
 "log"
 "net/http"
 "os"

 _ "modernc.org/sqlite" // Драйвер для SQLite

 "expenses/internal/handlers"
 "expenses/internal/repository" // Наш новый пакет
)

func main() {
 // 1. Подключаем базу данных SQLite
 db, err := sql.Open("sqlite", "expenses.db")
 if err != nil {
  log.Fatal("Ошибка открытия БД:", err)
 }
 defer db.Close()

// 2. Применяем миграцию (читаем SQL файл и создаем таблицу)
migration, err := os.ReadFile("migrations/001_init.sql")
if err != nil {
   log.Fatal("Ошибка чтения файла миграции:", err)
}
   
log.Println("Миграция прочитана, применяем...")
_, err = db.Exec(string(migration))
if err != nil {
   log.Fatal("Ошибка применения миграции:", err)
}
log.Println("Миграция успешно применена!")

 // 3. Инициализируем репозиторий (слой работы с БД)
 repo := repository.NewExpenseRepository(db)

 // 4. Шаблоны и обработчики
 tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
 // Передаем репозиторий в хендлер (вам нужно будет обновить и его структуру, см. ниже)
 h := handlers.NewExpenseHandler(tmpl, repo)

 // 5. Маршруты
 mux := http.NewServeMux()

 mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
   http.NotFound(w, r)
   return
  }
  http.Redirect(w, r, "/expenses", http.StatusSeeOther)
 })

 // Используем обработчики из h (убедитесь, что в handlers.go они привязаны к структуре)
 mux.HandleFunc("/expenses", func(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
   h.List(w, r)
  case http.MethodPost:
   h.Create(w, r)
  default:
   w.WriteHeader(http.StatusMethodNotAllowed)
  }
 })

 // Форма создания (явно указываем GET)
    mux.HandleFunc("GET /expenses/new", h.NewForm)

    // Форма редактирования
    mux.HandleFunc("GET /expenses/{id}/edit", h.EditForm)

    // Сохранение изменений (добавили /edit в конце)
    mux.HandleFunc("POST /expenses/{id}/edit", h.Update)

    // Удаление
    mux.HandleFunc("POST /expenses/{id}/delete", h.Delete)

 // Статика
 fs := http.FileServer(http.Dir("web/static"))
 mux.Handle("/static/", http.StripPrefix("/static/", fs))

 // Выносим порт в переменную
 port := ":8080"

 // Выводим полные ссылки
 log.Println("Сервер успешно запущен!")
 log.Printf("Главная страница: http://localhost%s/expenses", port)
 log.Printf("Форма добавления: http://localhost%s/expenses/new", port)

 // Запускаем сервер
 log.Fatal(http.ListenAndServe(port, mux))
}