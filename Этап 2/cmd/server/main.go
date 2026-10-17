package main

import (
 "html/template"
 "log"
 "net/http"

 "expenses/internal/handlers"
)
func main() {
 tmpl := template.Must(template.ParseGlob("web/templates/*.html"))

 h := handlers.NewExpenseHandler(tmpl)
 mux := http.NewServeMux()
 mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
   http.NotFound(w, r)
   return
  }
  http.Redirect(w, r, "/expenses", http.StatusSeeOther)
 })
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
 mux.HandleFunc("/expenses/new", h.NewForm)

 fs := http.FileServer(http.Dir("web/static"))
 mux.Handle("/static/", http.StripPrefix("/static/", fs))

 port := ":8080"

 log.Printf("Сервер успешно запущен!")
 log.Printf("Главная страница: http://localhost%s/expenses", port)
 log.Printf("Форма добавления: http://localhost%s/expenses/new", port)

 log.Fatal(http.ListenAndServe(port, mux))
}
