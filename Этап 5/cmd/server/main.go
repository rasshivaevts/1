package main

import (
   "database/sql"
   "html/template"
   "log"
   "net/http"
   "os"

   _ "modernc.org/sqlite"

   "expenses/internal/handlers"
   "expenses/internal/repository"
)

func main() {
   // 1. Подключаем базу данных SQLite
   db, err := sql.Open("sqlite", "expenses.db")
   if err != nil {
      log.Fatal("Ошибка открытия БД:", err)
   }
   defer db.Close()

   // 2. Применяем миграцию
   migration, err := os.ReadFile("migrations/001_init.sql")
   if err != nil {
      log.Fatal("Ошибка чтения файла миграции:", err)
   }
   if _, err := db.Exec(string(migration)); err != nil {
      log.Fatal("Ошибка применения миграции:", err)
   }

   // 3. Инициализируем репозитории
   expenseRepo := repository.NewExpenseRepository(db)
   userRepo := repository.NewUserRepository(db)

   // 4. Инициализируем обработчики и шаблоны
   tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
   
   expH := handlers.NewExpenseHandler(tmpl, expenseRepo)
   authH := handlers.NewAuthHandler(tmpl, userRepo)

   // 5. Настраиваем маршруты
   mux := http.NewServeMux()

   // --- Открытые маршруты (Авторизация) ---
   mux.HandleFunc("GET /login", authH.LoginPage)
   mux.HandleFunc("POST /login", authH.Login)
   mux.HandleFunc("GET /register", authH.RegisterPage)
   mux.HandleFunc("POST /register", authH.Register)
   mux.HandleFunc("GET /logout", authH.Logout)

   // --- Защищенные маршруты (Траты) ---
   // Оборачиваем в authH.RequireAuth
   mux.HandleFunc("/expenses", authH.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
      switch r.Method {
      case http.MethodGet:
         expH.List(w, r)
      case http.MethodPost:
         expH.Create(w, r)
      default:
         w.WriteHeader(http.StatusMethodNotAllowed)
      }
   }))

   mux.HandleFunc("GET /expenses/new", authH.RequireAuth(expH.NewForm))
   mux.HandleFunc("GET /expenses/{id}/edit", authH.RequireAuth(expH.EditForm))
   mux.HandleFunc("POST /expenses/{id}/edit", authH.RequireAuth(expH.Update))
   mux.HandleFunc("POST /expenses/{id}/delete", authH.RequireAuth(expH.Delete))

   // Главная страница (редирект)
   mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      if r.URL.Path != "/" {
         http.NotFound(w, r)
         return
      }
      http.Redirect(w, r, "/expenses", http.StatusSeeOther)
   })

   // Статика
   fs := http.FileServer(http.Dir("web/static"))
   mux.Handle("/static/", http.StripPrefix("/static/", fs))

   // 6. Запуск сервера
   port := ":8080"
   log.Println("Сервер успешно запущен!")
   log.Printf("Главная страница: http://localhost%s/expenses", port)
   log.Fatal(http.ListenAndServe(port, mux))
}