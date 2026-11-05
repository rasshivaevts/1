package models

import "time"

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type Expense struct {
	ID           int
	Amount       float64
	Description  string
	Date         time.Time
	CategoryID   int
	CategoryName string
	UserID       int // Привязка к пользователю
}

type Category struct {
	ID   int
	Name string
}