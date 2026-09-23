package main

import (
	"log"
	"net/http"

	"expenses_tracker/internal/handlers"
	"expenses_tracker/internal/models"
)

func main() {
	// Инициализируем хранилище трат в памяти
	store := models.NewExpenseStore()

	// Инициализируем обработчики с хранилищем и шаблонами
	h, err := handlers.NewHandlers(store)
	if err != nil {
		log.Fatalf("Не удалось инициализировать обработчики: %v", err)
	}

	// Регистрация маршрутов (этап 1 + этап 2)
	mux := http.NewServeMux()

	// --- Этап 1: базовые маршруты ---
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/about", h.About)
	mux.HandleFunc("/ping", h.Ping)

	// --- Этап 2: маршруты трат ---
	mux.HandleFunc("GET /expenses", h.ListExpenses)
	mux.HandleFunc("GET /expenses/new", h.NewExpense)
	mux.HandleFunc("POST /expenses", h.CreateExpense)

	// Раздача статических файлов (CSS)
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
