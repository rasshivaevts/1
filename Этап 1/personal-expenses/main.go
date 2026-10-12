package main

import (
	"fmt"
	"net/http"
)

// 1. Главная страница (текст-приветствие)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "Добро пожаловать в приложение учёта расходов!")
}

// 2. Страница /about (описание проекта)
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "Это учебный проект по веб-разработке на Go")
}

// 3. Страница /ping (проверка работы сервера и методов)
func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed) // 405
		fmt.Fprint(w, "Метод не разрешен")
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "pong")
}

func main() {
 mux := http.NewServeMux()

 mux.HandleFunc("/", homeHandler)
 mux.HandleFunc("/about", aboutHandler)
 mux.HandleFunc("/ping", pingHandler)

 port := ":8080"
 fmt.Println("Сервер успешно запущен!")
 fmt.Println("Доступные страницы:")
 fmt.Printf("  Главная:  http://localhost%s/\n", port)
 fmt.Printf("  О проекте: http://localhost%s/about\n", port)
 fmt.Printf("  Пинг:     http://localhost%s/ping\n", port)
 fmt.Println("Для остановки нажмите Ctrl+C")

 err := http.ListenAndServe(port, mux)
 if err != nil {
  fmt.Printf("Ошибка запуска сервера: %v\n", err)
 }
}