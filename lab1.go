package main

import (
 "fmt"
 "log"
 "net/http"
)

// Обработчик главной страницы "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
 // В Go http.ServeMux сопоставляет "/" со всеми ненайденными путями,
 // поэтому делаем точную проверку URL.Path
 if r.URL.Path != "/" {
  http.NotFound(w, r)
  return
 }
 w.Header().Set("Content-Type", "text/plain; charset=utf-8")
 w.WriteHeader(http.StatusOK)
 fmt.Fprintln(w, "Добро пожаловать на главный сервер!")
}

// Обработчик страницы "/about"
func aboutHandler(w http.ResponseWriter, r *http.Request) {
 if r.URL.Path != "/about" {
  http.NotFound(w, r)
  return
 }
 w.Header().Set("Content-Type", "text/plain; charset=utf-8")
 w.WriteHeader(http.StatusOK)
 fmt.Fprintln(w, "Описание проекта: Простой HTTP-сервер на Go без сторонних библиотек.")
}

// Обработчик маршрута "/ping" с проверкой метода (Задание 3)
func pingHandler(w http.ResponseWriter, r *http.Request) {
 if r.URL.Path != "/ping" {
  http.NotFound(w, r)
  return
 }

 // Разрешаем только метод GET, на все остальные возвращаем 405 Method Not Allowed
 if r.Method != http.MethodGet {
  w.Header().Set("Allow", http.MethodGet)
  http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed) // Статус 405
  return
 }

 w.Header().Set("Content-Type", "text/plain; charset=utf-8")
 w.WriteHeader(http.StatusOK) // Статус 200
 fmt.Fprintln(w, "pong")
}

func main() {
 mux := http.NewServeMux()

 // Регистрация маршрутов
 mux.HandleFunc("/", homeHandler)
 mux.HandleFunc("/about", aboutHandler)
 mux.HandleFunc("/ping", pingHandler)

 fmt.Println("Сервер запущен на http://localhost:8080")
 if err := http.ListenAndServe(":8080", mux); err != nil {
  log.Fatalf("Ошибка запуска сервера: %v", err)
 }
}